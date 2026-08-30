package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	dataTypeAttrName string = "DataType"
	deletedAttrName  string = "Deleted"
	// Each period is roughly 3.5 days
	periodDurationSecs  int64 = HistoryExpirationIntervalSecs / 4
	CurrentCountVersion int   = 2
)

// ClientItemCounts is used to marshal and unmarshal ClientItemCounts items in
// dynamoDB.
type ClientItemCounts struct {
	ClientID                string
	ID                      string
	ItemCount               int
	HistoryItemCountPeriod1 int
	HistoryItemCountPeriod2 int
	HistoryItemCountPeriod3 int
	HistoryItemCountPeriod4 int
	LastPeriodChangeTime    int64
	Version                 int
}

// ClientItemCountByClientID  implements sort.Interface for []ClientItemCount
// based on ClientID.
type ClientItemCountByClientID []ClientItemCounts

func (a ClientItemCountByClientID) Len() int      { return len(a) }
func (a ClientItemCountByClientID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ClientItemCountByClientID) Less(i, j int) bool {
	return a[i].ClientID < a[j].ClientID
}

func (counts *ClientItemCounts) SumHistoryCounts() int {
	return counts.HistoryItemCountPeriod1 +
		counts.HistoryItemCountPeriod2 +
		counts.HistoryItemCountPeriod3 +
		counts.HistoryItemCountPeriod4
}

func (dynamo *Dynamo) initRealCountsAndUpdateHistoryCounts(ctx context.Context, counts *ClientItemCounts) error {
	now := time.Now().Unix()
	if counts.Version < CurrentCountVersion {
		return dynamo.migrateLegacyItemCounts(ctx, counts, now)
	}
	rotateHistoryCountPeriods(counts, now)
	return nil
}

func (dynamo *Dynamo) migrateLegacyItemCounts(ctx context.Context, counts *ClientItemCounts, now int64) error {
	if counts.ItemCount > 0 {
		if err := dynamo.refreshItemCountsFromTable(ctx, counts); err != nil {
			return err
		}
	}
	counts.LastPeriodChangeTime = now
	counts.Version = CurrentCountVersion
	return nil
}

func (dynamo *Dynamo) countMatchingItems(
	ctx context.Context,
	pkCond expression.KeyConditionBuilder,
	filterCond expression.ConditionBuilder,
	kind string,
) (int, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(pkCond).WithFilter(filterCond).Build()
	if err != nil {
		return 0, fmt.Errorf("error building %s item count query: %w", kind, err)
	}
	out, err := dynamo.Query(ctx, &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(Table),
		Select:                    types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("error querying %s item count: %w", kind, err)
	}
	return int(out.Count), nil
}

func (dynamo *Dynamo) refreshItemCountsFromTable(ctx context.Context, counts *ClientItemCounts) error {
	// If last period change time is 0, assume that the old count exists in
	// ItemCount, which may include history items that have expired.
	pkCond := expression.Key(clientIDDataTypeMtimeIdxPk).Equal(expression.Value(counts.ClientID))
	historyFilter := expression.And(
		expression.Name(dataTypeAttrName).
			In(expression.Value(HistoryTypeID), expression.Value(HistoryDeleteDirectiveTypeID)),
		expression.Name(deletedAttrName).Equal(expression.Value(false)),
	)
	historyCount, err := dynamo.countMatchingItems(ctx, pkCond, historyFilter, "history")
	if err != nil {
		return err
	}
	counts.HistoryItemCountPeriod1 = 0
	counts.HistoryItemCountPeriod2 = 0
	counts.HistoryItemCountPeriod3 = 0
	counts.HistoryItemCountPeriod4 = historyCount

	normalFilter := expression.And(
		expression.AttributeExists(expression.Name(dataTypeAttrName)),
		expression.Name(dataTypeAttrName).NotEqual(expression.Value(HistoryTypeID)),
		expression.Name(dataTypeAttrName).NotEqual(expression.Value(HistoryDeleteDirectiveTypeID)),
		expression.Name(deletedAttrName).Equal(expression.Value(false)),
	)
	normalCount, err := dynamo.countMatchingItems(ctx, pkCond, normalFilter, "normal")
	if err != nil {
		return err
	}
	counts.ItemCount = normalCount
	return nil
}

func rotateHistoryCountPeriods(counts *ClientItemCounts, now int64) {
	timeSinceLastChange := now - counts.LastPeriodChangeTime
	if timeSinceLastChange < periodDurationSecs {
		return
	}
	changeCount := int(timeSinceLastChange / periodDurationSecs)
	for range changeCount {
		// The records from "period 1"/the earliest period
		// will be purged from the count, since they will be deleted via DDB TTL
		counts.HistoryItemCountPeriod1 = counts.HistoryItemCountPeriod2
		counts.HistoryItemCountPeriod2 = counts.HistoryItemCountPeriod3
		counts.HistoryItemCountPeriod3 = counts.HistoryItemCountPeriod4
		counts.HistoryItemCountPeriod4 = 0
	}
	counts.LastPeriodChangeTime += periodDurationSecs * int64(changeCount)
}

// GetClientItemCount returns the count of non-deleted sync items stored for
// a given client.
func (dynamo *Dynamo) GetClientItemCount(ctx context.Context, clientID string) (*ClientItemCounts, error) {
	primaryKey := PrimaryKey{ClientID: clientID, ID: clientID}
	key, err := attributevalue.MarshalMap(primaryKey)
	if err != nil {
		return nil, fmt.Errorf("error marshalling primary key to get item-count item: %w", err)
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(Table),
	}

	out, err := dynamo.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting an item-count item: %w", err)
	}

	clientItemCounts := &ClientItemCounts{}
	err = attributevalue.UnmarshalMap(out.Item, clientItemCounts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling item-count item: %w", err)
	}

	if len(clientItemCounts.ClientID) == 0 {
		clientItemCounts.ClientID = clientID
		clientItemCounts.ID = clientID
	}

	if err = dynamo.initRealCountsAndUpdateHistoryCounts(ctx, clientItemCounts); err != nil {
		return nil, err
	}

	return clientItemCounts, nil
}

// UpdateClientItemCount updates the count of non-deleted sync items for a
// given client stored in the dynamoDB.
func (dynamo *Dynamo) UpdateClientItemCount(
	ctx context.Context,
	counts *ClientItemCounts,
	newNormalItemCount int,
	newHistoryItemCount int,
) error {
	counts.HistoryItemCountPeriod4 += newHistoryItemCount
	counts.ItemCount += newNormalItemCount

	item, err := attributevalue.MarshalMap(*counts)
	if err != nil {
		return fmt.Errorf("error marshalling item counts: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(Table),
	}

	_, err = dynamo.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("error updating item-count item in dynamoDB: %w", err)
	}
	return nil
}
