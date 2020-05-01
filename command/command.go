package command

import (
	"encoding/binary"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/sync_pb"
	"github.com/rs/zerolog/log"
)

var (
	maxGUBatchSize int32 = 500 // Could be modified in tests.
)

const (
	storeBirthday              string = "1"
	maxCommitBatchSize         int32  = 90
	sessionsCommitDelaySeconds int32  = 11
	setSyncPollInterval        int32  = 30
	nigoriTypeID               int32  = 47745
)

// handleGetUpdatesRequest handles GetUpdatesMessage and fills
// GetUpdatesResponse. Target sync entities in the database will be updated or
// deleted based on the client's requests.
func handleGetUpdatesRequest(guMsg *sync_pb.GetUpdatesMessage, guRsp *sync_pb.GetUpdatesResponse, db datastore.Datastore, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later

	if *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT {
		err := InsertServerDefinedUniqueEntities(db, clientID)
		if err != nil {
			log.Error().Err(err).Msg("Create server defined unique entities failed")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, nil
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

		if *fromProgressMarker.DataTypeId == nigoriTypeID &&
			*guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT {
			// Bypassing chromium's restriction here, our server won't provide the
			// initial encryption keys like chromium does, this will be overwritten
			// by our client.
			guRsp.EncryptionKeys = make([][]byte, 1)
			guRsp.EncryptionKeys[0] = []byte("1234")
		}

		token, n := binary.Varint(guRsp.NewProgressMarker[i].Token)
		if n <= 0 {
			return &errCode, fmt.Errorf("Failed at decoding token value %v", token)
		}

		curMaxSize := int64(maxSize) - int64(len(guRsp.Entries))
		count, entities, err := db.GetUpdatesForType(int(*fromProgressMarker.DataTypeId), token, fetchFolders, clientID, curMaxSize)
		if err != nil {
			log.Error().Err(err).Msg("db.GetUpdatesForType failed")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, nil
		}

		// Fill the PB entry from above DB entries until maxSize is reached.
		j := 0
		for ; j < len(entities) && len(guRsp.Entries) < cap(guRsp.Entries); j++ {
			entity, err := datastore.CreatePBSyncEntity(&entities[j])
			if err != nil {
				errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
				return &errCode, nil
			}
			guRsp.Entries = append(guRsp.Entries, entity)
		}
		// Add to changesRemaining if this type has some items left due to batchSize.
		changesRemaining = count - int64(len(entities))

		// If entities are appended, use the lastest mtime as returned token.
		if j != 0 {
			guRsp.NewProgressMarker[i].Token = make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(guRsp.NewProgressMarker[i].Token, *entities[j-1].Mtime)
		}
	}

	return &errCode, nil
}

// handleCommitRequest handles the commit message and fills the commit response.
// New sync entity is created and inserted into the database.
func handleCommitRequest(commitMsg *sync_pb.CommitMessage, commitRsp *sync_pb.CommitResponse, db datastore.Datastore, clientID string) (*sync_pb.SyncEnums_ErrorType, error) {
	if commitMsg == nil {
		return nil, fmt.Errorf("nil commitMsg is received")
	}

	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	commitRsp.Entryresponse = make([]*sync_pb.CommitResponse_EntryResponse, len(commitMsg.Entries))
	if commitMsg.Entries != nil {
		for i, v := range commitMsg.Entries {
			entryRsp := &sync_pb.CommitResponse_EntryResponse{}
			commitRsp.Entryresponse[i] = entryRsp
			entityToCommit, err := datastore.CreateDBSyncEntity(v, commitMsg.CacheGuid, clientID)
			if err != nil { // Can't unmarshal & marshal the message from PB into DB format
				rspType := sync_pb.CommitResponse_INVALID_MESSAGE
				entryRsp.ResponseType = &rspType
				continue
			}

			*entityToCommit.Version++
			if *entityToCommit.Version == 1 { // Create
				err = db.InsertSyncEntity(entityToCommit)
				if err != nil {
					log.Error().Err(err).Msg("Insert sync entity failed")
					rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
					entryRsp.ResponseType = &rspType
					continue
				}
			} else { // Update
				conflict, err := db.UpdateSyncEntity(entityToCommit)
				if err != nil {
					log.Error().Err(err).Msg("Update sync entity failed")
					rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
					entryRsp.ResponseType = &rspType
					continue
				}
				if conflict {
					rspType := sync_pb.CommitResponse_CONFLICT
					entryRsp.ResponseType = &rspType
					continue
				}
			}

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
	}
	return &errCode, nil
}

// HandleClientToServerMessage handles the protobuf ClientToServerMessage and
// fills the protobuf ClientToServerResponse.
func HandleClientToServerMessage(pb *sync_pb.ClientToServerMessage, pbRsp *sync_pb.ClientToServerResponse, db datastore.Datastore, clientID string) error {
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
		pbRsp.ErrorCode, err = handleGetUpdatesRequest(pb.GetUpdates, guRsp, db, clientID)
	} else if *pb.MessageContents == sync_pb.ClientToServerMessage_COMMIT {
		commitRsp := &sync_pb.CommitResponse{}
		pbRsp.Commit = commitRsp
		pbRsp.ErrorCode, err = handleCommitRequest(pb.Commit, commitRsp, db, clientID)
	} else {
		return fmt.Errorf("unsupported message type of ClientToServerMessage")
	}

	return err
}
