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
	maxGUBatchSize       int32 = 500
	maxClientObjectQuota int   = 50000
)

const (
	storeBirthday              string = "1"
	maxCommitBatchSize         int32  = 90
	sessionsCommitDelaySeconds int32  = 110
	setSyncPollInterval        int32  = 1200
	nigoriTypeID               int32  = 47745
	deviceInfoTypeID           int    = 154522
	maxActiveDevices           int    = 50
)

// handleGetUpdatesRequest handles GetUpdatesMessage and fills
// GetUpdatesResponse. Target sync entities in the database will be updated or
// deleted based on the client's requests.
func handleGetUpdatesRequest(cache *cache.Cache, guMsg *sync_pb.GetUpdatesMessage, guRsp *sync_pb.GetUpdatesResponse, db datastore.Datastore, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	isNewClient := guMsg.GetUpdatesOrigin != nil && *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT
	isPoll := guMsg.GetUpdatesOrigin != nil && *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_PERIODIC
	if isNewClient {
		// Reject the request if client has >= 50 devices in the chain.
		activeDevices := 0
		for {
			hasChangesRemaining, syncEntities, err := db.GetUpdatesForType(deviceInfoTypeID, 0, false, clientID, int64(maxGUBatchSize))
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
		err := InsertServerDefinedUniqueEntities(db, clientID)
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
	if guMsg.BatchSize != nil && *guMsg.BatchSize < maxGUBatchSize {
		maxSize = *guMsg.BatchSize
	}

	// Process from_progress_marker
	guRsp.NewProgressMarker = make([]*sync_pb.DataTypeProgressMarker, len(guMsg.FromProgressMarker))
	guRsp.Entries = make([]*sync_pb.SyncEntity, 0, maxSize)
	for i, fromProgressMarker := range guMsg.FromProgressMarker {
		guRsp.NewProgressMarker[i] = &sync_pb.DataTypeProgressMarker{}
		guRsp.NewProgressMarker[i].DataTypeId = fromProgressMarker.DataTypeId

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
		if int32(len(guRsp.Entries)) >= maxSize {
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

		curMaxSize := int64(maxSize) - int64(len(guRsp.Entries))
		hasChangesRemaining, entities, err := db.GetUpdatesForType(int(*fromProgressMarker.DataTypeId), token, fetchFolders, clientID, curMaxSize)
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
			token == 0 && len(entities) == 0 {
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, fmt.Errorf("nigori root folder entity is not ready yet")
		}

		if hasChangesRemaining {
			changesRemaining = 1 // Chromium uses 1 instead of actual count of update entries remaining.
		}

		// Fill the PB entry from above DB entries until maxSize is reached.
		j := 0
		for ; j < len(entities) && len(guRsp.Entries) < cap(guRsp.Entries); j++ {
			entity, err := datastore.CreatePBSyncEntity(&entities[j])
			if err != nil {
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode, fmt.Errorf("error creating protobuf sync entity from DB entity: %w", err)
			}
			guRsp.Entries = append(guRsp.Entries, entity)
		}
		// If entities are appended, use the lastest mtime as returned token.
		if j != 0 {
			guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(guRsp.NewProgressMarker[i].Token, *entities[j-1].Mtime)
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
				mtime = *entities[j-1].Mtime
			}
			cache.SetTypeMtime(context.Background(), clientID, int(*fromProgressMarker.DataTypeId), mtime)
		}
	}

	return &errCode, nil
}

// handleCommitRequest handles the commit message and fills the commit response.
// For each commit entry:
//   - new sync entity is created and inserted into the database if version is 0.
//   - existed sync entity will be updated if version is greater than 0.
func handleCommitRequest(cache *cache.Cache, commitMsg *sync_pb.CommitMessage, commitRsp *sync_pb.CommitResponse, db datastore.Datastore, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	if commitMsg == nil {
		return nil, fmt.Errorf("nil commitMsg is received")
	}

	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	if commitMsg.Entries == nil {        // nothing to process
		return &errCode, nil
	}

	itemCount, err := db.GetClientItemCount(clientID)
	count := 0
	if err != nil {
		log.Error().Err(err).Msg("Get client's item count failed")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, fmt.Errorf("error getting client's item count: %w", err)
	}

	commitRsp.Entryresponse = make([]*sync_pb.CommitResponse_EntryResponse, len(commitMsg.Entries))

	// Map client-generated ID to its server-generated ID.
	idMap := make(map[string]string)
	// Map to save commit data type ID & mtime
	typeMtimeMap := make(map[int]int64)
	for i, v := range commitMsg.Entries {
		entryRsp := &sync_pb.CommitResponse_EntryResponse{}
		commitRsp.Entryresponse[i] = entryRsp

		entityToCommit, err := datastore.CreateDBSyncEntity(v, commitMsg.CacheGuid, clientID)
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
		*entityToCommit.Version = *entityToCommit.Mtime
		if oldVersion == 0 { // Create
			if itemCount+count >= maxClientObjectQuota {
				rspType := sync_pb.CommitResponse_OVER_QUOTA
				entryRsp.ResponseType = &rspType
				entryRsp.ErrorMessage = aws.String(fmt.Sprintf("There are already %v non-deleted objects in store", itemCount))
				continue
			}

			conflict, err := db.InsertSyncEntity(entityToCommit)
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

			count++
		} else { // Update
			conflict, delete, err := db.UpdateSyncEntity(entityToCommit, oldVersion)
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
			if delete {
				count--
			}
		}

		typeMtimeMap[*entityToCommit.DataType] = *entityToCommit.Mtime
		// Prepare success response
		rspType := sync_pb.CommitResponse_SUCCESS
		entryRsp.ResponseType = &rspType
		entryRsp.IdString = aws.String(entityToCommit.ID)
		entryRsp.Version = entityToCommit.Version
		entryRsp.ParentIdString = entityToCommit.ParentID
		entryRsp.Name = entityToCommit.Name
		entryRsp.NonUniqueName = entityToCommit.NonUniqueName
		entryRsp.Mtime = entityToCommit.Mtime
	}

	// Save (clientID#dataType, mtime) into cache after writing into DB.
	for dataType, mtime := range typeMtimeMap {
		cache.SetTypeMtime(context.Background(), clientID, dataType, mtime)
	}

	err = db.UpdateClientItemCount(clientID, count)
	if err != nil {
		// We only impose a soft quota limit on the item count for each client, so
		// we only log the error without further actions here. The reason of this
		// is we do not want to pay the cost to ensure strong consistency on this
		// value and we do not want to give up previous DB operations if we cannot
		// update the count this time. In addition, we do not retry this operation
		// either because it is acceptable to miss one time of this update and
		// chances of failing to update the item count multiple times in a row for
		// a single client is quite low.
		log.Error().Err(err).Msg("Update client item count failed")
	}
	return &errCode, nil
}

// HandleClientToServerMessage handles the protobuf ClientToServerMessage and
// fills the protobuf ClientToServerResponse.
func HandleClientToServerMessage(cache *cache.Cache, pb *sync_pb.ClientToServerMessage, pbRsp *sync_pb.ClientToServerResponse, db datastore.Datastore, clientID string) error {
	// Create ClientToServerResponse and fill general fields for both GU and
	// Commit.
	pbRsp.StoreBirthday = aws.String(storeBirthday)
	pbRsp.ClientCommand = &sync_pb.ClientCommand{
		SetSyncPollInterval:        aws.Int32(setSyncPollInterval),
		MaxCommitBatchSize:         aws.Int32(maxCommitBatchSize),
		SessionsCommitDelaySeconds: aws.Int32(sessionsCommitDelaySeconds)}

	var err error
	if pb.MessageContents == nil {
		return fmt.Errorf("nil pb.MessageContents received")
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_GET_UPDATES {
		guRsp := &sync_pb.GetUpdatesResponse{}
		pbRsp.GetUpdates = guRsp
		pbRsp.ErrorCode, err = handleGetUpdatesRequest(cache, pb.GetUpdates, guRsp, db, clientID)
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
		pbRsp.ErrorCode, err = handleCommitRequest(cache, pb.Commit, commitRsp, db, clientID)
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
	} else {
		return fmt.Errorf("unsupported message type of ClientToServerMessage")
	}

	return nil
}
