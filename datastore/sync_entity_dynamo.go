package datastore

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/brave/go-sync/utils"
)

const (
	maxBatchGetItemSize       = 100 // Limited by AWS.
	maxTransactDeleteItemSize = 10  // Limited by AWS.
	clientTagItemPrefix       = "Client#"
	serverTagItemPrefix       = "Server#"
	conditionalCheckFailed    = "ConditionalCheckFailed"
	disabledChainID           = "disabled_chain"
	reasonDeleted             = "deleted"
)

// SyncEntityByClientIDID implements sort.Interface for []SyncEntity based on
// the string concatenation of ClientID and ID fields.
type SyncEntityByClientIDID []SyncEntity

func (a SyncEntityByClientIDID) Len() int      { return len(a) }
func (a SyncEntityByClientIDID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SyncEntityByClientIDID) Less(i, j int) bool {
	return a[i].ClientID+a[i].ID < a[j].ClientID+a[j].ID
}

// SyncEntityByMtime implements sort.Interface for []SyncEntity based on Mtime.
type SyncEntityByMtime []SyncEntity

func (a SyncEntityByMtime) Len() int      { return len(a) }
func (a SyncEntityByMtime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SyncEntityByMtime) Less(i, j int) bool {
	return *a[i].Mtime < *a[j].Mtime
}

// DisabledMarkerItem is used to mark sync chain as deleted in Dynamodb
type DisabledMarkerItem struct {
	ClientID string
	ID       string
	Reason   string
	Mtime    *int64
	Ctime    *int64
}

// DisabledMarkerItemQuery is used to query for disabled marker item in
// DynamoDB
type DisabledMarkerItemQuery struct {
	ClientID string
	ID       string
}

// ServerClientUniqueTagItem is used to marshal and unmarshal tag items in
// dynamoDB.
type ServerClientUniqueTagItem struct {
	ClientID string // Hash key
	ID       string // Range key
	Mtime    *int64
	Ctime    *int64
}

// ServerClientUniqueTagItemQuery is used to query for unique tag items in
// dynamoDB.
type ServerClientUniqueTagItemQuery struct {
	ClientID string // Hash key
	ID       string // Range key
}

// TagItemByClientIDID implements sort.Interface for []ServerClientUniqueTagItem
// based on the string concatenation of ClientID and ID fields.
type TagItemByClientIDID []ServerClientUniqueTagItem

func (a TagItemByClientIDID) Len() int      { return len(a) }
func (a TagItemByClientIDID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a TagItemByClientIDID) Less(i, j int) bool {
	return a[i].ClientID+a[i].ID < a[j].ClientID+a[j].ID
}

// getTagPrefix is a helper method to give the proper prefix for unique tag
func getTagPrefix(isServer bool) string {
	if isServer {
		return serverTagItemPrefix
	}
	return clientTagItemPrefix
}

// NewServerClientUniqueTagItem creates a tag item which is used to ensure the
// uniqueness of server-defined or client-defined unique tags for a client.
func NewServerClientUniqueTagItem(clientID string, tag string, isServer bool) *ServerClientUniqueTagItem {
	prefix := getTagPrefix(isServer)
	now := aws.Int64(utils.UnixMilli(time.Now()))

	return &ServerClientUniqueTagItem{
		ClientID: clientID,
		ID:       prefix + tag,
		Mtime:    now,
		Ctime:    now,
	}
}

// NewServerClientUniqueTagItemQuery creates a tag item query which is used to
// determine whether a sync entity has a unique tag item or not
func NewServerClientUniqueTagItemQuery(clientID string, tag string, isServer bool) *ServerClientUniqueTagItemQuery {
	prefix := getTagPrefix(isServer)

	return &ServerClientUniqueTagItemQuery{
		ClientID: clientID,
		ID:       prefix + tag,
	}
}

// InsertSyncEntity inserts a new sync entity into dynamoDB.
// If ClientDefinedUniqueTag is not null, we will use a write transaction to
// write a sync item along with a tag item to ensure the uniqueness of the
// client tag. Otherwise, only a sync item is written into DB without using
// transactions.
func (dynamo *Dynamo) InsertSyncEntity(entity *SyncEntity) (bool, error) {
	// Create a condition for inserting new items only.
	cond := expression.AttributeNotExists(expression.Name(pk))
	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return false, fmt.Errorf("error building expression to insert sync entity: %w", err)
	}

	// Write tag item for all data types, except for
	// the history type, which does not use tag items.
	if entity.ClientDefinedUniqueTag != nil && *entity.DataType != HistoryTypeID {
		items := []*dynamodb.TransactWriteItem{}
		// Additional item for ensuring tag's uniqueness for a specific client.
		item := NewServerClientUniqueTagItem(entity.ClientID, *entity.ClientDefinedUniqueTag, false)
		av, err := dynamodbattribute.MarshalMap(*item)
		if err != nil {
			return false, fmt.Errorf("error marshalling unique tag item to insert sync entity: %w", err)
		}
		tagItem := &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:                      av,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				ConditionExpression:       expr.Condition(),
				TableName:                 aws.String(Table),
			},
		}

		// Normal sync item
		av, err = dynamodbattribute.MarshalMap(*entity)
		if err != nil {
			return false, fmt.Errorf("error marshlling sync item to insert sync entity: %w", err)
		}
		syncItem := &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:                      av,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				ConditionExpression:       expr.Condition(),
				TableName:                 aws.String(Table),
			},
		}
		items = append(items, tagItem)
		items = append(items, syncItem)

		_, err = dynamo.TransactWriteItems(
			&dynamodb.TransactWriteItemsInput{TransactItems: items})
		if err != nil {
			// Return conflict if insert condition failed.
			if canceledException, ok := err.(*dynamodb.TransactionCanceledException); ok {
				for _, reason := range canceledException.CancellationReasons {
					if reason.Code != nil && *reason.Code == conditionalCheckFailed {
						return true, fmt.Errorf("error inserting sync item with client tag: %w", err)
					}
				}
			}
			return false, fmt.Errorf("error writing tag item and sync item in a transaction to insert sync entity: %w", err)
		}

		return false, nil
	}

	// Normal sync item
	av, err := dynamodbattribute.MarshalMap(*entity)
	if err != nil {
		return false, fmt.Errorf("error marshalling sync item to insert sync entity: %w", err)
	}
	input := &dynamodb.PutItemInput{
		Item:                      av,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		TableName:                 aws.String(Table),
	}
	_, err = dynamo.PutItem(input)
	if err != nil {
		return false, fmt.Errorf("error calling PutItem to insert sync item: %w", err)
	}
	return false, nil
}

// HasServerDefinedUniqueTag check the tag item to see if there is already a
// tag item exists with the tag value for a specific client.
func (dynamo *Dynamo) HasServerDefinedUniqueTag(clientID string, tag string) (bool, error) {
	tagItem := NewServerClientUniqueTagItemQuery(clientID, tag, true)
	key, err := dynamodbattribute.MarshalMap(tagItem)
	if err != nil {
		return false, fmt.Errorf("error marshalling key to check if server tag existed: %w", err)
	}

	input := &dynamodb.GetItemInput{
		Key:                  key,
		ProjectionExpression: aws.String(projPk),
		TableName:            aws.String(Table),
	}

	out, err := dynamo.GetItem(input)
	if err != nil {
		return false, fmt.Errorf("error calling GetItem to check if server tag existed: %w", err)
	}

	return out.Item != nil, nil
}

func (dynamo *Dynamo) HasItem(clientID string, ID string) (bool, error) {
	primaryKey := PrimaryKey{ClientID: clientID, ID: ID}
	key, err := dynamodbattribute.MarshalMap(primaryKey)

	if err != nil {
		return false, fmt.Errorf("error marshalling key to check if item existed: %w", err)
	}

	input := &dynamodb.GetItemInput{
		Key:                  key,
		ProjectionExpression: aws.String(projPk),
		TableName:            aws.String(Table),
	}

	out, err := dynamo.GetItem(input)
	if err != nil {
		return false, fmt.Errorf("error calling GetItem to check if item existed: %w", err)
	}

	return out.Item != nil, nil
}

// InsertSyncEntitiesWithServerTags is used to insert sync entities with
// server-defined unique tags. To ensure the uniqueness, for each sync entity,
// we will write a tag item and a sync item. Items for all the entities in the
// array would be written into DB in one transaction.
func (dynamo *Dynamo) InsertSyncEntitiesWithServerTags(entities []*SyncEntity) error {
	items := []*dynamodb.TransactWriteItem{}
	for _, entity := range entities {
		// Create a condition for inserting new items only.
		cond := expression.AttributeNotExists(expression.Name(pk))
		expr, err := expression.NewBuilder().WithCondition(cond).Build()
		if err != nil {
			return fmt.Errorf("error building expression to insert sync entity with server tag: %w", err)
		}

		// Additional item for ensuring tag's uniqueness for a specific client.
		item := NewServerClientUniqueTagItem(entity.ClientID, *entity.ServerDefinedUniqueTag, true)
		av, err := dynamodbattribute.MarshalMap(*item)
		if err != nil {
			return fmt.Errorf("error marshalling tag item to insert sync entity with server tag: %w", err)
		}
		tagItem := &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:                      av,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				ConditionExpression:       expr.Condition(),
				TableName:                 aws.String(Table),
			},
		}

		// Normal sync item
		av, err = dynamodbattribute.MarshalMap(*entity)
		if err != nil {
			return fmt.Errorf("error marshalling sync item to insert sync entity with server tag: %w", err)
		}
		syncItem := &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:                      av,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				ConditionExpression:       expr.Condition(),
				TableName:                 aws.String(Table),
			},
		}

		items = append(items, tagItem)
		items = append(items, syncItem)
	}

	_, err := dynamo.TransactWriteItems(
		&dynamodb.TransactWriteItemsInput{TransactItems: items})
	if err != nil {
		return fmt.Errorf("error writing sync entities with server tags in a transaction: %w", err)
	}
	return nil
}

// DisableSyncChain marks a chain as disabled so no further updates or commits can happen
func (dynamo *Dynamo) DisableSyncChain(clientID string) error {
	now := aws.Int64(utils.UnixMilli(time.Now()))
	disabledMarker := DisabledMarkerItem{
		ClientID: clientID,
		ID:       disabledChainID,
		Reason:   reasonDeleted,
		Mtime:    now,
		Ctime:    now,
	}

	av, err := dynamodbattribute.MarshalMap(disabledMarker)
	if err != nil {
		return fmt.Errorf("error marshalling disabled marker: %w", err)
	}

	markerInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(Table),
	}

	_, err = dynamo.PutItem(markerInput)
	if err != nil {
		return fmt.Errorf("error calling PutItem to insert sync item: %w", err)
	}

	return nil
}

// ClearServerData deletes all items for a given clientID
func (dynamo *Dynamo) ClearServerData(clientID string) ([]SyncEntity, error) {
	syncEntities := []SyncEntity{}
	pkb := expression.Key(pk)
	pkv := expression.Value(clientID)
	keyCond := expression.KeyEqual(pkb, pkv)
	exprs := expression.NewBuilder().WithKeyCondition(keyCond)
	expr, err := exprs.Build()
	if err != nil {
		return syncEntities, fmt.Errorf("error building expression to get updates: %w", err)
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(Table),
	}

	out, err := dynamo.Query(input)
	if err != nil {
		return syncEntities, fmt.Errorf("error doing query to get updates: %w", err)
	}
	count := *out.Count

	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &syncEntities)
	if err != nil {
		return syncEntities, fmt.Errorf("error unmarshalling updated sync entities: %w", err)
	}

	var i, j int64
	for i = 0; i < count; i += maxTransactDeleteItemSize {
		j = i + maxTransactDeleteItemSize
		if j > count {
			j = count
		}

		items := []*dynamodb.TransactWriteItem{}
		for _, item := range syncEntities[i:j] {
			if item.ID == disabledChainID {
				continue
			}

			// Fail delete if race condition detected (modified time has changed).
			if item.Version != nil {
				cond := expression.Name("Mtime").Equal(expression.Value(*item.Mtime))
				expr, err := expression.NewBuilder().WithCondition(cond).Build()
				if err != nil {
					return syncEntities, fmt.Errorf("error deleting sync entities for client %s: %w", clientID, err)
				}

				writeItem := dynamodb.TransactWriteItem{
					Delete: &dynamodb.Delete{
						ConditionExpression:       expr.Condition(),
						ExpressionAttributeNames:  expr.Names(),
						ExpressionAttributeValues: expr.Values(),
						TableName:                 aws.String(Table),
						Key: map[string]*dynamodb.AttributeValue{
							pk: {
								S: aws.String(item.ClientID),
							},
							sk: {
								S: aws.String(item.ID),
							},
						},
					},
				}

				items = append(items, &writeItem)
			} else {
				// If row doesn't hold Mtime, delete as usual.
				writeItem := dynamodb.TransactWriteItem{
					Delete: &dynamodb.Delete{
						TableName: aws.String(Table),
						Key: map[string]*dynamodb.AttributeValue{
							pk: {
								S: aws.String(item.ClientID),
							},
							sk: {
								S: aws.String(item.ID),
							},
						},
					},
				}

				items = append(items, &writeItem)
			}

		}

		_, err = dynamo.TransactWriteItems(&dynamodb.TransactWriteItemsInput{TransactItems: items})
		if err != nil {
			return syncEntities, fmt.Errorf("error deleting sync entities for client %s: %w", clientID, err)
		}
	}

	return syncEntities, nil
}

// IsSyncChainDisabled checks whether a given sync chain has been deleted
func (dynamo *Dynamo) IsSyncChainDisabled(clientID string) (bool, error) {
	key, err := dynamodbattribute.MarshalMap(DisabledMarkerItemQuery{
		ClientID: clientID,
		ID:       disabledChainID,
	})
	if err != nil {
		return false, fmt.Errorf("error marshalling key to check if server tag existed: %w", err)
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(Table),
	}

	out, err := dynamo.GetItem(input)
	if err != nil {
		return false, fmt.Errorf("error calling GetItem to check if sync chain disabled: %w", err)
	}

	return len(out.Item) > 0, nil
}

// UpdateSyncEntity updates a sync item in dynamoDB.
func (dynamo *Dynamo) UpdateSyncEntity(entity *SyncEntity, oldVersion int64) (bool, bool, error) {
	primaryKey := PrimaryKey{ClientID: entity.ClientID, ID: entity.ID}
	key, err := dynamodbattribute.MarshalMap(primaryKey)
	if err != nil {
		return false, false, fmt.Errorf("error marshalling key to update sync entity: %w", err)
	}

	// condition to ensure the request is update only...
	cond := expression.AttributeExists(expression.Name(pk))
	// ...and the version matches, if applicable
	if *entity.DataType != HistoryTypeID {
		cond = expression.And(cond, expression.Name("Version").Equal(expression.Value(oldVersion)))
	}

	update := expression.Set(expression.Name("Version"), expression.Value(entity.Version))
	update = update.Set(expression.Name("Mtime"), expression.Value(entity.Mtime))
	update = update.Set(expression.Name("Specifics"), expression.Value(entity.Specifics))
	update = update.Set(expression.Name("DataTypeMtime"), expression.Value(entity.DataTypeMtime))

	// Update optional fields only if the value is not null.
	if entity.UniquePosition != nil {
		update = update.Set(expression.Name("UniquePosition"), expression.Value(entity.UniquePosition))
	}
	if entity.ParentID != nil {
		update = update.Set(expression.Name("ParentID"), expression.Value(entity.ParentID))
	}
	if entity.Name != nil {
		update = update.Set(expression.Name("Name"), expression.Value(entity.Name))
	}
	if entity.NonUniqueName != nil {
		update = update.Set(expression.Name("NonUniqueName"), expression.Value(entity.NonUniqueName))
	}
	if entity.Deleted != nil {
		update = update.Set(expression.Name("Deleted"), expression.Value(entity.Deleted))
	}
	if entity.Folder != nil {
		update = update.Set(expression.Name("Folder"), expression.Value(entity.Folder))
	}

	expr, err := expression.NewBuilder().WithCondition(cond).WithUpdate(update).Build()
	if err != nil {
		return false, false, fmt.Errorf("error building expression to update sync entity: %w", err)
	}

	// Soft-delete a sync item with a client tag, use a transaction to delete its
	// tag item too.
	if entity.Deleted != nil && entity.ClientDefinedUniqueTag != nil && *entity.Deleted && *entity.DataType != HistoryTypeID {
		pk := PrimaryKey{
			ClientID: entity.ClientID, ID: clientTagItemPrefix + *entity.ClientDefinedUniqueTag}
		tagItemKey, err := dynamodbattribute.MarshalMap(pk)
		if err != nil {
			return false, false, fmt.Errorf("error marshalling key to update sync entity: %w", err)
		}

		items := []*dynamodb.TransactWriteItem{}
		updateSyncItem := &dynamodb.TransactWriteItem{
			Update: &dynamodb.Update{
				Key:                                 key,
				ExpressionAttributeNames:            expr.Names(),
				ExpressionAttributeValues:           expr.Values(),
				ConditionExpression:                 expr.Condition(),
				UpdateExpression:                    expr.Update(),
				ReturnValuesOnConditionCheckFailure: aws.String(dynamodb.ReturnValueAllOld),
				TableName:                           aws.String(Table),
			},
		}
		deleteTagItem := &dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				Key:       tagItemKey,
				TableName: aws.String(Table),
			},
		}
		items = append(items, updateSyncItem)
		items = append(items, deleteTagItem)

		_, err = dynamo.TransactWriteItems(
			&dynamodb.TransactWriteItemsInput{TransactItems: items})
		if err != nil {
			// Return conflict if the update condition fails.
			if canceledException, ok := err.(*dynamodb.TransactionCanceledException); ok {
				for _, reason := range canceledException.CancellationReasons {
					if reason.Code != nil && *reason.Code == conditionalCheckFailed {
						return true, false, nil
					}
				}
			}

			return false, false, fmt.Errorf("error deleting sync item and tag item in a transaction: %w", err)
		}

		// Successfully soft-delete the sync item and delete the tag item.
		return false, true, nil
	}

	// Not deleting a sync item with a client tag, do a normal update on sync
	// item.
	input := &dynamodb.UpdateItemInput{
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              aws.String(dynamodb.ReturnValueAllOld),
		TableName:                 aws.String(Table),
	}

	out, err := dynamo.UpdateItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			// Return conflict if the write condition fails.
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return true, false, nil
			}
		}
		return false, false, fmt.Errorf("error calling UpdateItem to update sync entity: %w", err)
	}

	// Unmarshal out.Attributes
	oldEntity := &SyncEntity{}
	err = dynamodbattribute.UnmarshalMap(out.Attributes, oldEntity)
	if err != nil {
		return false, false, fmt.Errorf("error unmarshalling old sync entity: %w", err)
	}
	var deleted bool
	if entity.Deleted == nil { // No updates on Deleted this time.
		deleted = false
	} else if oldEntity.Deleted == nil { // Consider it as Deleted = false.
		deleted = *entity.Deleted
	} else {
		deleted = !*oldEntity.Deleted && *entity.Deleted
	}
	return false, deleted, nil
}

// GetUpdatesForType returns sync entities of a data type where it's mtime is
// later than the client token.
// To do this in dynamoDB, we use (ClientID, DataType#Mtime) as GSI to get a
// list of (ClientID, ID) primary keys with the given condition, then read the
// actual sync item using the list of primary keys.
func (dynamo *Dynamo) GetUpdatesForType(dataType int, clientToken int64, fetchFolders bool, clientID string, maxSize int64) (bool, []SyncEntity, error) {
	syncEntities := []SyncEntity{}

	// Get (ClientID, ID) pairs which are updates after mtime for a data type,
	// sorted by dataType#mTime. e.g. sorted by mtime since dataType is the same.
	dataTypeMtimeLowerBound := strconv.Itoa(dataType) + "#" + strconv.FormatInt(clientToken+1, 10)
	dataTypeMtimeUpperBound := strconv.Itoa(dataType+1) + "#0"
	pkCond := expression.Key(clientIDDataTypeMtimeIdxPk).Equal(expression.Value(clientID))
	skCond := expression.KeyBetween(
		expression.Key(clientIDDataTypeMtimeIdxSk),
		expression.Value(dataTypeMtimeLowerBound),
		expression.Value(dataTypeMtimeUpperBound))
	keyCond := expression.KeyAnd(pkCond, skCond)
	exprs := expression.NewBuilder().WithKeyCondition(keyCond)

	if !fetchFolders { // Filter folder entities out if fetchFolder is false.
		exprs = exprs.WithFilter(
			expression.Equal(expression.Name("Folder"), expression.Value(false)))
	}

	expr, err := exprs.Build()
	if err != nil {
		return false, syncEntities, fmt.Errorf("error building expression to get updates: %w", err)
	}

	input := &dynamodb.QueryInput{
		IndexName:                 aws.String(clientIDDataTypeMtimeIdx),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      aws.String(projPk),
		TableName:                 aws.String(Table),
		Limit:                     aws.Int64(maxSize),
	}

	out, err := dynamo.Query(input)
	if err != nil {
		return false, syncEntities, fmt.Errorf("error doing query to get updates: %w", err)
	}

	hasChangesRemaining := false
	if out.LastEvaluatedKey != nil && len(out.LastEvaluatedKey) > 0 {
		hasChangesRemaining = true
	}

	count := *(out.Count)
	if count == 0 { // No updates
		return hasChangesRemaining, syncEntities, nil
	}

	// Use return (ClientID, ID) primary keys to get the actual items.
	var outAv []map[string]*dynamodb.AttributeValue
	var i, j int64
	for i = 0; i < count; i += maxBatchGetItemSize {
		j = i + maxBatchGetItemSize
		if j > count {
			j = count
		}

		batchInput := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				Table: {
					Keys: out.Items[i:j],
				},
			},
		}

		err := dynamo.BatchGetItemPages(batchInput,
			func(batchOut *dynamodb.BatchGetItemOutput, last bool) bool {
				outAv = append(outAv, batchOut.Responses[Table]...)
				return last
			})
		if err != nil {
			return false, syncEntities, fmt.Errorf("error getting update items in a batch: %w", err)
		}
	}

	err = dynamodbattribute.UnmarshalListOfMaps(outAv, &syncEntities)
	if err != nil {
		return false, syncEntities, fmt.Errorf("error unmarshalling updated sync entities: %w", err)
	}

	// filter out any expired items, i.e. history sync entities over 90 days old
	nowUnix := time.Now().Unix()
	var filteredSyncEntities []SyncEntity
	for _, syncEntity := range syncEntities {
		if syncEntity.ExpirationTime != nil && *syncEntity.ExpirationTime > 0 {
			if *syncEntity.ExpirationTime < nowUnix {
				continue
			}
		}
		filteredSyncEntities = append(filteredSyncEntities, syncEntity)
	}

	sort.Sort(SyncEntityByMtime(filteredSyncEntities))
	return hasChangesRemaining, filteredSyncEntities, nil
}
