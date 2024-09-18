package command

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/rs/zerolog/log"
)

var (
	// Could be modified in tests.
	maxGUBatchSize              = 500
	maxClientObjectQuota        = 50000
	maxClientHistoryObjectQuota = 30000
)

const (
	storeBirthday       string = "1"
	maxCommitBatchSize  int32  = 90
	setSyncPollInterval int32  = 30
	nigoriTypeID        int32  = 47745
	deviceInfoTypeID    int    = 154522
	maxActiveDevices    int    = 50
	historyCountTypeStr string = "history"
	normalCountTypeStr  string = "normal"
)

// handleGetUpdatesRequest handles GetUpdatesMessage and fills
// GetUpdatesResponse. Target sync entities in the database will be updated or
// deleted based on the client's requests.
func handleGetUpdatesRequest(cache *cache.Cache, guMsg *sync_pb.GetUpdatesMessage, guRsp *sync_pb.GetUpdatesResponse, dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDatastore, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	isNewClient := guMsg.GetUpdatesOrigin != nil && *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT
	isPoll := guMsg.GetUpdatesOrigin != nil && *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_PERIODIC

	dbHelpers, err := NewDBHelpers(dynamoDB, sqlDB, clientID, nil, false)
	if err != nil {
		return nil, err
	}
	defer dbHelpers.Trx.Rollback()

	if isNewClient {
		// Reject the request if client has >= 50 devices in the chain.
		activeDevices := 0
		for {
			hasChangesRemaining, syncEntities, err := dbHelpers.getUpdatesFromDBs(deviceInfoTypeID, 0, false, maxGUBatchSize)
			if err != nil {
				log.Error().Err(err).Msgf("db.GetUpdatesForType failed for type %v", deviceInfoTypeID)
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode,
					fmt.Errorf("error getting updates for type %v: %w", deviceInfoTypeID, err)
			}

			for _, entity := range syncEntities {
				if !*entity.Deleted {
					activeDevices++
				}

				// Error out when exceeds the limit.
				if activeDevices >= maxActiveDevices {
					errCode = sync_pb.SyncEnums_THROTTLED
					return &errCode, fmt.Errorf("exceed limit of active devices in a chain")
				}
			}

			// Run until all device records are checked.
			if !hasChangesRemaining {
				break
			}
		}

		// Insert initial records if needed.
		err := dbHelpers.InsertServerDefinedUniqueEntities()
		if err != nil {
			log.Error().Err(err).Msg("Create server defined unique entities failed")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, fmt.Errorf("error creating server defined unique entitiies: %w", err)
		}
	}

	changesRemaining := int64(0)
	guRsp.ChangesRemaining = &changesRemaining

	if guMsg.FromProgressMarker == nil { // nothing to process
		return &errCode, nil
	}

	fetchFolders := true
	if guMsg.FetchFolders != nil {
		fetchFolders = *guMsg.FetchFolders
	}

	maxSize := maxGUBatchSize

	// Process from_progress_marker
	guRsp.NewProgressMarker = make([]*sync_pb.DataTypeProgressMarker, len(guMsg.FromProgressMarker))
	guRsp.Entries = make([]*sync_pb.SyncEntity, 0, maxSize)

	var dataTypes []int

	for i, fromProgressMarker := range guMsg.FromProgressMarker {
		guRsp.NewProgressMarker[i] = &sync_pb.DataTypeProgressMarker{}
		guRsp.NewProgressMarker[i].DataTypeId = fromProgressMarker.DataTypeId

		dataTypes = append(dataTypes, int(*fromProgressMarker.DataTypeId))

		// Default token value is client's token, otherwise 0.
		// This token will be updated when we return the updated entities.
		if len(fromProgressMarker.Token) > 0 {
			guRsp.NewProgressMarker[i].Token = fromProgressMarker.Token
		} else {
			guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(guRsp.NewProgressMarker[i].Token, int64(0))
		}

		if *fromProgressMarker.DataTypeId == nigoriTypeID && isNewClient {
			// Bypassing chromium's restriction here, our server won't provide the
			// initial encryption keys like chromium does, this will be overwritten
			// by our client.
			guRsp.EncryptionKeys = make([][]byte, 1)
			guRsp.EncryptionKeys[0] = []byte("1234")
		}

		// No need to get updates for this type because we already reach the
		// maximum GetUpdates size for this request. Continue to next type instead
		// of break because we need to prepare NewProgressMarker for all entries in
		// FromProgressMarker, where the returned token stays the same as the one
		// passed in FromProgressMarker.
		if len(guRsp.Entries) >= maxSize {
			continue
		}

		token, n := binary.Varint(guRsp.NewProgressMarker[i].Token)
		if n <= 0 {
			return nil, fmt.Errorf("Failed at decoding token value %v", token)
		}

		// Check cache to short circuit with 0 updates for polling requests.
		if isPoll &&
			!cache.IsTypeMtimeUpdated(context.Background(), clientID, int(*fromProgressMarker.DataTypeId), token) {
			continue
		}

		curMaxSize := maxSize - len(guRsp.Entries)
		hasChangesRemaining, syncEntities, err := dbHelpers.getUpdatesFromDBs(int(*fromProgressMarker.DataTypeId), token, fetchFolders, curMaxSize)
		if err != nil {
			log.Error().Err(err).Msgf("db.GetUpdatesForType failed for type %v", *fromProgressMarker.DataTypeId)
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode,
				fmt.Errorf("error getting updates for type %v: %w", *fromProgressMarker.DataTypeId, err)
		}

		// Due to eventually read consistency, it is possible that we cannot get
		// the nigori root folder entity for this NEW_CLIENT GetUpdates request,
		// which is essential for clients when initializing sync engine with nigori
		// type. Return a transient error for clients to re-request in this case.
		if isNewClient && *fromProgressMarker.DataTypeId == nigoriTypeID &&
			token == 0 && len(syncEntities) == 0 {
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, fmt.Errorf("nigori root folder entity is not ready yet")
		}

		if hasChangesRemaining {
			changesRemaining = 1 // Chromium uses 1 instead of actual count of update entries remaining.
		}

		// Fill the PB entry from above DB entries until maxSize is reached.
		j := 0
		for ; j < len(syncEntities) && len(guRsp.Entries) < cap(guRsp.Entries); j++ {
			entity, err := datastore.CreatePBSyncEntity(&syncEntities[j])
			if err != nil {
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode, fmt.Errorf("error creating protobuf sync entity from DB entity: %w", err)
			}
			guRsp.Entries = append(guRsp.Entries, entity)
		}
		// If entities are appended, use the lastest mtime as returned token.
		if j != 0 {
			guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(guRsp.NewProgressMarker[i].Token, *syncEntities[j-1].Mtime)
		}

		// Save (clientID#dataType, mtime) into cache after querying from DB.
		// If changes_remaining = 1 in the response, client will send another poll
		// request immediately, we do not save mtime into cache in this iteration
		// because the client token in the subsequent poll request will be equal to
		// this mtime and we will wrongly think there are no updates when we
		// process that subsequent poll request. The cache will be updated in a
		// subsequent poll request where changes_remaining = 0.
		if changesRemaining != 1 {
			var mtime int64
			if j == 0 {
				mtime = token
			} else {
				mtime = *syncEntities[j-1].Mtime
			}
			cache.SetTypeMtime(context.Background(), clientID, int(*fromProgressMarker.DataTypeId), mtime)
		}
	}

	migratedEntities, err := dbHelpers.maybeMigrateToSQL(dataTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to perform migration: %w", err)
	}

	if len(migratedEntities) > 0 {
		if err = dynamoDB.DeleteEntities(migratedEntities); err != nil {
			log.Error().Err(err).Msgf("Failed to delete migrated items")
		}
	}

	if err = dbHelpers.Trx.Commit(); err != nil {
		return nil, err
	}

	return &errCode, nil
}

// handleCommitRequest handles the commit message and fills the commit response.
// For each commit entry:
//   - new sync entity is created and inserted into the database if version is 0.
//   - existed sync entity will be updated if version is greater than 0.
func handleCommitRequest(cache *cache.Cache, commitMsg *sync_pb.CommitMessage, commitRsp *sync_pb.CommitResponse, dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDatastore, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	if commitMsg == nil {
		return nil, fmt.Errorf("nil commitMsg is received")
	}

	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	if commitMsg.Entries == nil {        // nothing to process
		return &errCode, nil
	}

	if !sqlDB.Variations().Ready {
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, fmt.Errorf("SQL rollout not ready")
	}

	dbHelpers, err := NewDBHelpers(dynamoDB, sqlDB, clientID, cache, true)
	if err != nil {
		return nil, err
	}
	defer dbHelpers.Trx.Rollback()

	commitRsp.Entryresponse = make([]*sync_pb.CommitResponse_EntryResponse, len(commitMsg.Entries))

	// Map client-generated ID to its server-generated ID.
	idMap := make(map[string]string)
	// Map to save commit data type ID & mtime
	typeMtimeMap := make(map[int]int64)

	var migratedEntities []*datastore.SyncEntity
	for i, v := range commitMsg.Entries {
		entryRsp := &sync_pb.CommitResponse_EntryResponse{}
		commitRsp.Entryresponse[i] = entryRsp

		entityToCommit, err := datastore.CreateDBSyncEntity(v, commitMsg.CacheGuid, clientID, dbHelpers.ChainID)
		if err != nil { // Can't unmarshal & marshal the message from PB into DB format
			rspType := sync_pb.CommitResponse_INVALID_MESSAGE
			entryRsp.ResponseType = &rspType
			entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Cannot convert protobuf sync entity to DB format: %v", err.Error()))
			continue
		}

		// Check if ParentID is a client-generated ID which appears in previous
		// commit entries, if so, replace with corresponding server-generated ID.
		if entityToCommit.ParentID != nil {
			if serverParentID, ok := idMap[*entityToCommit.ParentID]; ok {
				entityToCommit.ParentID = &serverParentID
			}
		}

		oldVersion := *entityToCommit.Version
		isUpdateOp := oldVersion != 0
		isHistoryItem := *entityToCommit.DataType == datastore.HistoryTypeID
		isHistoryRelatedItem := isHistoryItem || *entityToCommit.DataType == datastore.HistoryDeleteDirectiveTypeID
		*entityToCommit.Version = *entityToCommit.Mtime

		if isHistoryItem {
			isUpdateOp, err = dbHelpers.hasItemInEitherDB(entityToCommit)
			if err != nil {
				log.Error().Err(err).Msg("Insert history sync entity failed")
				rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
				entryRsp.ResponseType = &rspType
				entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Insert history sync entity failed: %v", err.Error()))
				continue
			}
		}

		if !isUpdateOp { // Create
			totalItemCount := dbHelpers.ItemCounts.sumCounts(false)
			if totalItemCount >= maxClientObjectQuota {
				rspType := sync_pb.CommitResponse_OVER_QUOTA
				entryRsp.ResponseType = &rspType
				entryRsp.ErrorMessage = aws.String(fmt.Sprintf("There are already %v non-deleted objects in store", totalItemCount))
				continue
			}

			if !isHistoryRelatedItem || dbHelpers.ItemCounts.sumCounts(true) < maxClientHistoryObjectQuota {
				// Insert all non-history items. For history items, ignore any items above history quoto
				// and lie to the client about the objects being synced instead of returning OVER_QUOTA
				// so the client can continue to sync other entities.
				var conflict bool
				conflict, err = dbHelpers.insertSyncEntity(entityToCommit)
				if err != nil {
					log.Error().Err(err).Msg("Insert sync entity failed")
					rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
					if conflict {
						rspType = sync_pb.CommitResponse_CONFLICT
					}
					entryRsp.ResponseType = &rspType
					entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Insert sync entity failed: %v", err.Error()))
					continue
				}

				// Save client-generated to server-generated ID mapping when committing
				// a new entry with OriginatorClientItemID (client-generated ID).
				if entityToCommit.OriginatorClientItemID != nil {
					idMap[*entityToCommit.OriginatorClientItemID] = entityToCommit.ID
				}
			}
		} else { // Update
			conflict, migratedEntity, err := dbHelpers.updateSyncEntity(entityToCommit, oldVersion)
			if err != nil {
				log.Error().Err(err).Msg("Update sync entity failed")
				rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
				entryRsp.ResponseType = &rspType
				entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Update sync entity failed: %v", err.Error()))
				continue
			}
			if conflict {
				rspType := sync_pb.CommitResponse_CONFLICT
				entryRsp.ResponseType = &rspType
				continue
			}
			if migratedEntity != nil {
				migratedEntities = append(migratedEntities, migratedEntity)
			}
		}
		if err != nil {
			log.Error().Err(err).Msg("Interim count update failed")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, fmt.Errorf("Interim count update failed: %w", err)
		}

		typeMtimeMap[*entityToCommit.DataType] = *entityToCommit.Mtime
		// Prepare success response
		rspType := sync_pb.CommitResponse_SUCCESS
		entryRsp.ResponseType = &rspType
		entryRsp.IdString = aws.String(entityToCommit.ID)
		entryRsp.Version = entityToCommit.Version
		entryRsp.Mtime = entityToCommit.Mtime
	}

	err = dbHelpers.ItemCounts.save()
	if err != nil {
		log.Error().Err(err).Msg("Get interim item counts failed")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, fmt.Errorf("error getting interim item count: %w", err)
	}

	// Save (clientID#dataType, mtime) into cache after writing into DB.
	for dataType, mtime := range typeMtimeMap {
		cache.SetTypeMtime(context.Background(), clientID, dataType, mtime)
	}

	if len(migratedEntities) > 0 {
		if err = dynamoDB.DeleteEntities(migratedEntities); err != nil {
			log.Error().Err(err).Msgf("Failed to delete migrated items")
		}
	}

	if err = dbHelpers.Trx.Commit(); err != nil {
		return nil, err
	}

	return &errCode, nil
}

// handleClearServerDataRequest handles clearing user data from the datastore and cache
// and fills the response
func handleClearServerDataRequest(cache *cache.Cache, dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDatastore, _ *sync_pb.ClearServerDataMessage, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	errCode := sync_pb.SyncEnums_SUCCESS
	var err error

	dbHelpers, err := NewDBHelpers(dynamoDB, sqlDB, clientID, nil, false)
	if err != nil {
		return nil, err
	}
	defer dbHelpers.Trx.Rollback()

	err = dynamoDB.DisableSyncChain(clientID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to disable sync chain")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, err
	}

	syncEntities, err := dynamoDB.ClearServerData(clientID)
	if err != nil {
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, err
	}

	typeMtimeCacheKeys := []string{}
	for _, entity := range syncEntities {
		if entity.DataType != nil {
			typeMtimeCacheKeys = append(typeMtimeCacheKeys, cache.GetTypeMtimeKey(entity.ClientID, *entity.DataType))
		}
	}

	if len(typeMtimeCacheKeys) > 0 {
		err = cache.Del(context.Background(), typeMtimeCacheKeys...)
		if err != nil {
			log.Error().Err(err).Msg("Failed to clear cache")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, err
		}
	}

	if err = dbHelpers.SqlDB.DeleteChain(dbHelpers.Trx, dbHelpers.ChainID); err != nil {
		log.Error().Err(err).Msg("Failed to disable sync chain")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, err
	}

	if err = dbHelpers.Trx.Commit(); err != nil {
		return nil, err
	}

	return &errCode, nil
}

// HandleClientToServerMessage handles the protobuf ClientToServerMessage and
// fills the protobuf ClientToServerResponse.
func HandleClientToServerMessage(cache *cache.Cache, pb *sync_pb.ClientToServerMessage, pbRsp *sync_pb.ClientToServerResponse, dynamoDB datastore.DynamoDatastore, sqlDB datastore.SQLDatastore, clientID string) error {
	// Create ClientToServerResponse and fill general fields for both GU and
	// Commit.
	pbRsp.StoreBirthday = aws.String(storeBirthday)
	pbRsp.ClientCommand = &sync_pb.ClientCommand{
		SetSyncPollInterval: aws.Int32(setSyncPollInterval),
		MaxCommitBatchSize:  aws.Int32(maxCommitBatchSize)}

	var err error
	if pb.MessageContents == nil {
		return fmt.Errorf("nil pb.MessageContents received")
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_GET_UPDATES {
		guRsp := &sync_pb.GetUpdatesResponse{}
		pbRsp.GetUpdates = guRsp
		pbRsp.ErrorCode, err = handleGetUpdatesRequest(cache, pb.GetUpdates, guRsp, dynamoDB, sqlDB, clientID)
		if err != nil {
			if pbRsp.ErrorCode != nil {
				pbRsp.ErrorMessage = aws.String(err.Error())
				return nil
			}
			// In seledom error cases which are not temporary and will not go away
			// when clients retry, we will not use defined sync error in the proto
			// response, but use internal server error.
			return fmt.Errorf("error handling GetUpdates request: %w", err)
		}
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_COMMIT {
		commitRsp := &sync_pb.CommitResponse{}
		pbRsp.Commit = commitRsp
		pbRsp.ErrorCode, err = handleCommitRequest(cache, pb.Commit, commitRsp, dynamoDB, sqlDB, clientID)
		if err != nil {
			if pbRsp.ErrorCode != nil {
				pbRsp.ErrorMessage = aws.String(err.Error())
				return nil
			}
			// In seledom error cases which are not temporary and will not go away
			// when clients retry, we will not use defined sync error in the proto
			// response, but use internal server error.
			return fmt.Errorf("error handling Commit request: %w", err)
		}
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_CLEAR_SERVER_DATA {
		csdRsp := &sync_pb.ClearServerDataResponse{}
		pbRsp.ClearServerData = csdRsp
		pbRsp.ErrorCode, err = handleClearServerDataRequest(cache, dynamoDB, sqlDB, pb.ClearServerData, clientID)
		if err != nil {
			if pbRsp.ErrorCode != nil {
				pbRsp.ErrorMessage = aws.String(err.Error())
				return nil
			}
			// In seldom error cases which are not temporary and will not go away
			// when clients retry, we will not use defined sync error in the proto
			// response, but use internal server error.
			return fmt.Errorf("error handling ClearServerData request: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported message type of ClientToServerMessage")
	}

	return nil
}
