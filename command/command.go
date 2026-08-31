package command

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
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
	historyCountTypeStr string = "history"
	normalCountTypeStr  string = "normal"
)

func setupNewClient(
	ctx context.Context,
	db datastore.Datastore,
	clientID string,
) (*sync_pb.SyncEnums_ErrorType, error) {
	// Reject the request if client has >= 50 devices in the chain.
	activeDevices := 0
	for {
		hasChangesRemaining, syncEntities, err := db.GetUpdatesForType(
			ctx,
			deviceInfoTypeID,
			0,
			false,
			clientID,
			maxGUBatchSize,
		)
		if err != nil {
			log.Error().Err(err).Msgf("db.GetUpdatesForType failed for type %v", deviceInfoTypeID)
			errCode := sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode,
				fmt.Errorf("error getting updates for type %v: %w", deviceInfoTypeID, err)
		}

		for _, entity := range syncEntities {
			if !*entity.Deleted {
				activeDevices++
			}

			// Error out when device limit has been reached.
			if hasReachedDeviceLimit(activeDevices, clientID) {
				errCode := sync_pb.SyncEnums_THROTTLED
				return &errCode, errors.New("exceed limit of active devices in a chain")
			}
		}

		// Run until all device records are checked.
		if !hasChangesRemaining {
			break
		}
	}

	// Insert initial records if needed.
	err := InsertServerDefinedUniqueEntities(ctx, db, clientID)
	if err != nil {
		log.Error().Err(err).Msg("Create server defined unique entities failed")
		errCode := sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, fmt.Errorf("error creating server defined unique entities: %w", err)
	}
	success := sync_pb.SyncEnums_SUCCESS
	return &success, nil
}

func isNewClientOrigin(guMsg *sync_pb.GetUpdatesMessage) bool {
	return guMsg.GetUpdatesOrigin != nil && *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_NEW_CLIENT
}

func isPollOrigin(guMsg *sync_pb.GetUpdatesMessage) bool {
	return guMsg.GetUpdatesOrigin != nil && *guMsg.GetUpdatesOrigin == sync_pb.SyncEnums_PERIODIC
}

func maybeSetupNewClient(
	ctx context.Context,
	db datastore.Datastore,
	clientID string,
	isNewClient bool,
) (*sync_pb.SyncEnums_ErrorType, error) {
	if !isNewClient {
		return nil, nil //nolint:nilnil // no proto error code when setup is skipped
	}
	return setupNewClient(ctx, db, clientID)
}

func shouldFetchFolders(guMsg *sync_pb.GetUpdatesMessage) bool {
	if guMsg.FetchFolders != nil {
		return *guMsg.FetchFolders
	}
	return true
}

func progressMarkerToken(fromProgressMarker *sync_pb.DataTypeProgressMarker) []byte {
	// Default token value is client's token, otherwise 0.
	// This token will be updated when we return the updated entities.
	if len(fromProgressMarker.Token) > 0 {
		return fromProgressMarker.Token
	}
	token := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(token, int64(0))
	return token
}

func maybeSetNigoriEncryptionKeys(guRsp *sync_pb.GetUpdatesResponse, dataTypeID int32, isNewClient bool) {
	if dataTypeID != nigoriTypeID || !isNewClient {
		return
	}
	// Bypassing chromium's restriction here, our server won't provide the
	// initial encryption keys like chromium does, this will be overwritten
	// by our client.
	guRsp.EncryptionKeys = make([][]byte, 1)
	guRsp.EncryptionKeys[0] = []byte("1234")
}

func nigoriRootNotReady(isNewClient bool, dataTypeID int32, token int64, entities []datastore.SyncEntity) bool {
	// Due to eventually read consistency, it is possible that we cannot get
	// the nigori root folder entity for this NEW_CLIENT GetUpdates request,
	// which is essential for clients when initializing sync engine with nigori
	// type. Return a transient error for clients to re-request in this case.
	return isNewClient && dataTypeID == nigoriTypeID && token == 0 && len(entities) == 0
}

func appendGetUpdatesEntities(
	guRsp *sync_pb.GetUpdatesResponse,
	markerIndex int,
	entities []datastore.SyncEntity,
) (int, *sync_pb.SyncEnums_ErrorType, error) {
	// Fill the PB entry from above DB entries until maxSize is reached.
	j := 0
	for ; j < len(entities) && len(guRsp.Entries) < cap(guRsp.Entries); j++ {
		entity, createErr := datastore.CreatePBSyncEntity(&entities[j])
		if createErr != nil {
			errCode := sync_pb.SyncEnums_TRANSIENT_ERROR
			return 0, &errCode, fmt.Errorf("error creating protobuf sync entity from DB entity: %w", createErr)
		}
		guRsp.Entries = append(guRsp.Entries, entity)
	}
	// If entities are appended, use the lastest mtime as returned token.
	if j != 0 {
		guRsp.NewProgressMarker[markerIndex].Token = make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(guRsp.NewProgressMarker[markerIndex].Token, *entities[j-1].Mtime)
	}
	return j, nil, nil
}

func maybeSetTypeMtime(
	ctx context.Context,
	cache *cache.Cache,
	clientID string,
	dataTypeID int,
	changesRemaining int64,
	appendedCount int,
	token int64,
	entities []datastore.SyncEntity,
) {
	// Save (clientID#dataType, mtime) into cache after querying from DB.
	// If changes_remaining = 1 in the response, client will send another poll
	// request immediately, we do not save mtime into cache in this iteration
	// because the client token in the subsequent poll request will be equal to
	// this mtime and we will wrongly think there are no updates when we
	// process that subsequent poll request. The cache will be updated in a
	// subsequent poll request where changes_remaining = 0.
	if changesRemaining == 1 {
		return
	}
	mtime := token
	if appendedCount != 0 {
		mtime = *entities[appendedCount-1].Mtime
	}
	cache.SetTypeMtime(ctx, clientID, dataTypeID, mtime)
}

func processGetUpdatesForType(
	ctx context.Context,
	cache *cache.Cache,
	db datastore.Datastore,
	clientID string,
	fromProgressMarker *sync_pb.DataTypeProgressMarker,
	guRsp *sync_pb.GetUpdatesResponse,
	markerIndex int,
	maxSize int,
	fetchFolders bool,
	isNewClient bool,
	isPoll bool,
	changesRemaining *int64,
) (*sync_pb.SyncEnums_ErrorType, error) {
	guRsp.NewProgressMarker[markerIndex] = &sync_pb.DataTypeProgressMarker{
		DataTypeId: fromProgressMarker.DataTypeId,
		Token:      progressMarkerToken(fromProgressMarker),
	}

	maybeSetNigoriEncryptionKeys(guRsp, *fromProgressMarker.DataTypeId, isNewClient)

	// No need to get updates for this type because we already reach the
	// maximum GetUpdates size for this request. Continue to next type instead
	// of break because we need to prepare NewProgressMarker for all entries in
	// FromProgressMarker, where the returned token stays the same as the one
	// passed in FromProgressMarker.
	if len(guRsp.Entries) >= maxSize {
		return nil, nil //nolint:nilnil // continue to next type; no proto error
	}

	token, n := binary.Varint(guRsp.NewProgressMarker[markerIndex].Token)
	if n <= 0 {
		return nil, fmt.Errorf("failed at decoding token value %v", token)
	}

	// Check cache to short circuit with 0 updates for polling requests.
	if isPoll &&
		!cache.IsTypeMtimeUpdated(ctx, clientID, int(*fromProgressMarker.DataTypeId), token) {
		return nil, nil //nolint:nilnil // cache hit; no proto error
	}

	curMaxSize := maxSize - len(guRsp.Entries)
	hasChangesRemaining, entities, err := db.GetUpdatesForType(
		ctx,
		int(*fromProgressMarker.DataTypeId),
		token,
		fetchFolders,
		clientID,
		curMaxSize,
	)
	if err != nil {
		log.Error().Err(err).Msgf("db.GetUpdatesForType failed for type %v", *fromProgressMarker.DataTypeId)
		errCode := sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode,
			fmt.Errorf("error getting updates for type %v: %w", *fromProgressMarker.DataTypeId, err)
	}

	if nigoriRootNotReady(isNewClient, *fromProgressMarker.DataTypeId, token, entities) {
		errCode := sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, errors.New("nigori root folder entity is not ready yet")
	}

	if hasChangesRemaining {
		*changesRemaining = 1 // Chromium uses 1 instead of actual count of update entries remaining.
	}

	j, errCode, err := appendGetUpdatesEntities(guRsp, markerIndex, entities)
	if err != nil {
		return errCode, err
	}

	maybeSetTypeMtime(
		ctx,
		cache,
		clientID,
		int(*fromProgressMarker.DataTypeId),
		*changesRemaining,
		j,
		token,
		entities,
	)
	return nil, nil //nolint:nilnil // type processed; no proto error
}

// handleGetUpdatesRequest handles GetUpdatesMessage and fills
// GetUpdatesResponse. Target sync entities in the database will be updated or
// deleted based on the client's requests.
func handleGetUpdatesRequest(
	ctx context.Context,
	cache *cache.Cache,
	guMsg *sync_pb.GetUpdatesMessage,
	guRsp *sync_pb.GetUpdatesResponse,
	db datastore.Datastore,
	clientID string,
) (*sync_pb.SyncEnums_ErrorType, error) {
	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	isNewClient := isNewClientOrigin(guMsg)
	isPoll := isPollOrigin(guMsg)
	if setupErrCode, err := maybeSetupNewClient(ctx, db, clientID, isNewClient); err != nil {
		return setupErrCode, err
	}

	changesRemaining := int64(0)
	guRsp.ChangesRemaining = &changesRemaining

	if guMsg.FromProgressMarker == nil { // nothing to process
		return &errCode, nil
	}

	fetchFolders := shouldFetchFolders(guMsg)
	maxSize := maxGUBatchSize

	// Process from_progress_marker
	guRsp.NewProgressMarker = make([]*sync_pb.DataTypeProgressMarker, len(guMsg.FromProgressMarker))
	guRsp.Entries = make([]*sync_pb.SyncEntity, 0, maxSize)
	for i, fromProgressMarker := range guMsg.FromProgressMarker {
		typeErrCode, err := processGetUpdatesForType(
			ctx,
			cache,
			db,
			clientID,
			fromProgressMarker,
			guRsp,
			i,
			maxSize,
			fetchFolders,
			isNewClient,
			isPoll,
			&changesRemaining,
		)
		if err != nil {
			return typeErrCode, err
		}
	}

	return &errCode, nil
}

func getItemCounts(
	ctx context.Context,
	cache *cache.Cache,
	db datastore.Datastore,
	clientID string,
) (*datastore.ClientItemCounts, int, int, error) {
	itemCounts, err := db.GetClientItemCount(ctx, clientID)
	if err != nil {
		return nil, 0, 0, err
	}
	newNormalCount, newHistoryCount, err := getInterimItemCounts(ctx, cache, clientID, false)
	if err != nil {
		return nil, 0, 0, err
	}
	return itemCounts, newNormalCount, newHistoryCount, nil
}

func getInterimItemCounts(ctx context.Context, cache *cache.Cache, clientID string, clearCache bool) (int, int, error) {
	newNormalCount, err := cache.GetInterimCount(ctx, clientID, normalCountTypeStr, clearCache)
	if err != nil {
		return 0, 0, err
	}
	newHistoryCount, err := cache.GetInterimCount(ctx, clientID, historyCountTypeStr, clearCache)
	if err != nil {
		return 0, 0, err
	}
	return newNormalCount, newHistoryCount, nil
}

type commitCountState struct {
	currentNormal  int
	currentHistory int
	newNormal      int
	newHistory     int
	boostedQuota   int
}

func bumpInterimCount(
	ctx context.Context,
	cache *cache.Cache,
	clientID string,
	isHistoryRelatedItem bool,
	subtract bool,
	counts *commitCountState,
) error {
	countType := normalCountTypeStr
	if isHistoryRelatedItem {
		countType = historyCountTypeStr
	}
	newCount, err := cache.IncrementInterimCount(ctx, clientID, countType, subtract)
	if isHistoryRelatedItem {
		counts.newHistory = newCount
	} else {
		counts.newNormal = newCount
	}
	return err
}

func insertCommitEntity(
	ctx context.Context,
	db datastore.Datastore,
	cache *cache.Cache,
	clientID string,
	entity *datastore.SyncEntity,
	entryRsp *sync_pb.CommitResponse_EntryResponse,
	idMap map[string]string,
	counts *commitCountState,
	isHistoryRelatedItem bool,
) (bool, error) {
	total := counts.currentNormal + counts.currentHistory + counts.newNormal + counts.newHistory
	if total >= maxClientObjectQuota+counts.boostedQuota {
		rspType := sync_pb.CommitResponse_OVER_QUOTA
		entryRsp.ResponseType = &rspType
		entryRsp.ErrorMessage = aws.String(
			fmt.Sprintf(
				"There are already %v non-deleted objects in store",
				counts.currentNormal+counts.currentHistory,
			),
		)
		return true, nil
	}

	// Insert all non-history items. For history items, ignore any items above history quota
	// and lie to the client about the objects being synced instead of returning OVER_QUOTA
	// so the client can continue to sync other entities.
	if isHistoryRelatedItem && counts.currentHistory+counts.newHistory >= maxClientHistoryObjectQuota {
		return false, nil
	}

	conflict, insertErr := db.InsertSyncEntity(ctx, entity)
	if insertErr != nil {
		log.Error().Err(insertErr).Msg("Insert sync entity failed")
		rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
		if conflict {
			rspType = sync_pb.CommitResponse_CONFLICT
		}
		entryRsp.ResponseType = &rspType
		entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Insert sync entity failed: %v", insertErr.Error()))
		return true, nil
	}

	// Save client-generated to server-generated ID mapping when committing
	// a new entry with OriginatorClientItemID (client-generated ID).
	if entity.OriginatorClientItemID != nil {
		idMap[*entity.OriginatorClientItemID] = entity.ID
	}

	return false, bumpInterimCount(ctx, cache, clientID, isHistoryRelatedItem, false, counts)
}

func updateCommitEntity(
	ctx context.Context,
	db datastore.Datastore,
	cache *cache.Cache,
	clientID string,
	entity *datastore.SyncEntity,
	oldVersion int64,
	entryRsp *sync_pb.CommitResponse_EntryResponse,
	counts *commitCountState,
	isHistoryRelatedItem bool,
) (bool, error) {
	conflict, deleted, updateErr := db.UpdateSyncEntity(ctx, entity, oldVersion)
	if updateErr != nil {
		log.Error().Err(updateErr).Msg("Update sync entity failed")
		rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
		entryRsp.ResponseType = &rspType
		entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Update sync entity failed: %v", updateErr.Error()))
		return true, nil
	}
	if conflict {
		rspType := sync_pb.CommitResponse_CONFLICT
		entryRsp.ResponseType = &rspType
		return true, nil
	}
	if !deleted {
		return false, nil
	}
	return false, bumpInterimCount(ctx, cache, clientID, isHistoryRelatedItem, true, counts)
}

func applyBoostedHistoryQuota(counts *commitCountState) {
	if counts.currentHistory <= maxClientHistoryObjectQuota {
		return
	}
	// Sync chains with history entities stored before the history count fix
	// may have history counts greater than the new history item quota.
	// "Boost" the quota with the difference between the history quota and count,
	// so users can start syncing other entities immediately, instead of waiting for the
	// history TTL to get rid of the excess items.
	counts.boostedQuota = min(
		maxClientObjectQuota-maxClientHistoryObjectQuota,
		counts.currentHistory-maxClientHistoryObjectQuota,
	)
}

func replaceClientGeneratedParentID(entity *datastore.SyncEntity, idMap map[string]string) {
	// Check if ParentID is a client-generated ID which appears in previous
	// commit entries, if so, replace with corresponding server-generated ID.
	if entity.ParentID == nil {
		return
	}
	serverParentID, ok := idMap[*entity.ParentID]
	if !ok {
		return
	}
	entity.ParentID = &serverParentID
}

func resolveHistoryUpdateOp(
	ctx context.Context,
	db datastore.Datastore,
	clientID string,
	entity *datastore.SyncEntity,
	entryRsp *sync_pb.CommitResponse_EntryResponse,
) (bool, bool) {
	// Check if item exists using client_unique_tag
	isUpdateOp, err := db.HasItem(ctx, clientID, *entity.ClientDefinedUniqueTag)
	if err != nil {
		log.Error().Err(err).Msg("Insert history sync entity failed")
		rspType := sync_pb.CommitResponse_TRANSIENT_ERROR
		entryRsp.ResponseType = &rspType
		entryRsp.ErrorMessage = aws.String(fmt.Sprintf("Insert history sync entity failed: %v", err.Error()))
		return false, true
	}
	return isUpdateOp, false
}

func commitEntry(
	ctx context.Context,
	cache *cache.Cache,
	db datastore.Datastore,
	clientID string,
	pbEntity *sync_pb.SyncEntity,
	cacheGUID *string,
	entryRsp *sync_pb.CommitResponse_EntryResponse,
	idMap map[string]string,
	typeMtimeMap map[int]int64,
	counts *commitCountState,
	sizeMonitor *ItemSizeMonitor,
) error {
	entityToCommit, createErr := datastore.CreateDBSyncEntity(pbEntity, cacheGUID, clientID)
	if createErr != nil { // Can't unmarshal & marshal the message from PB into DB format
		rspType := sync_pb.CommitResponse_INVALID_MESSAGE
		entryRsp.ResponseType = &rspType
		entryRsp.ErrorMessage = aws.String(
			fmt.Sprintf("Cannot convert protobuf sync entity to DB format: %v", createErr.Error()),
		)
		return nil
	}

	sizeMonitor.Observe(*entityToCommit.DataType, pbEntity)
	replaceClientGeneratedParentID(entityToCommit, idMap)

	oldVersion := *entityToCommit.Version
	isUpdateOp := oldVersion != 0
	isHistoryRelatedItem := *entityToCommit.DataType == datastore.HistoryTypeID ||
		*entityToCommit.DataType == datastore.HistoryDeleteDirectiveTypeID
	*entityToCommit.Version = *entityToCommit.Mtime
	if *entityToCommit.DataType == datastore.HistoryTypeID {
		var skip bool
		isUpdateOp, skip = resolveHistoryUpdateOp(ctx, db, clientID, entityToCommit, entryRsp)
		if skip {
			return nil
		}
	}

	var skipEntry bool
	var err error
	if isUpdateOp {
		skipEntry, err = updateCommitEntity(
			ctx, db, cache, clientID, entityToCommit, oldVersion, entryRsp, counts, isHistoryRelatedItem,
		)
	} else {
		skipEntry, err = insertCommitEntity(
			ctx, db, cache, clientID, entityToCommit, entryRsp, idMap, counts, isHistoryRelatedItem,
		)
	}
	if err != nil {
		return err
	}
	if skipEntry {
		return nil
	}

	typeMtimeMap[*entityToCommit.DataType] = *entityToCommit.Mtime
	rspType := sync_pb.CommitResponse_SUCCESS
	entryRsp.ResponseType = &rspType
	entryRsp.IdString = aws.String(entityToCommit.ID)
	entryRsp.Version = entityToCommit.Version
	entryRsp.Mtime = entityToCommit.Mtime
	return nil
}

// handleCommitRequest handles the commit message and fills the commit response.
// For each commit entry:
//   - new sync entity is created and inserted into the database if version is 0.
//   - existed sync entity will be updated if version is greater than 0.
func handleCommitRequest(
	ctx context.Context,
	cache *cache.Cache,
	commitMsg *sync_pb.CommitMessage,
	commitRsp *sync_pb.CommitResponse,
	db datastore.Datastore,
	clientID string,
) (*sync_pb.SyncEnums_ErrorType, error) {
	if commitMsg == nil {
		return nil, errors.New("nil commitMsg is received")
	}

	errCode := sync_pb.SyncEnums_SUCCESS // default value, might be changed later
	if commitMsg.Entries == nil {        // nothing to process
		return &errCode, nil
	}

	itemCounts, newNormalCount, newHistoryCount, err := getItemCounts(ctx, cache, db, clientID)
	if err != nil {
		log.Error().Err(err).Msg("Get client's item count failed")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, fmt.Errorf("error getting client's item count: %w", err)
	}
	counts := commitCountState{
		currentNormal:  itemCounts.ItemCount,
		currentHistory: itemCounts.SumHistoryCounts(),
		newNormal:      newNormalCount,
		newHistory:     newHistoryCount,
	}
	applyBoostedHistoryQuota(&counts)

	commitRsp.Entryresponse = make([]*sync_pb.CommitResponse_EntryResponse, len(commitMsg.Entries))

	// Map client-generated ID to its server-generated ID.
	idMap := make(map[string]string)
	// Map to save commit data type ID & mtime
	typeMtimeMap := make(map[int]int64)
	sizeMonitor := NewItemSizeMonitor()
	for i, v := range commitMsg.Entries {
		entryRsp := &sync_pb.CommitResponse_EntryResponse{}
		commitRsp.Entryresponse[i] = entryRsp
		commitErr := commitEntry(
			ctx, cache, db, clientID, v, commitMsg.CacheGuid, entryRsp, idMap, typeMtimeMap, &counts, sizeMonitor,
		)
		if commitErr != nil {
			log.Error().Err(commitErr).Msg("Interim count update failed")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, fmt.Errorf("interim count update failed: %w", commitErr)
		}
	}

	sizeMonitor.LogWarnings()

	newNormalCount, newHistoryCount, err = getInterimItemCounts(ctx, cache, clientID, true)
	if err != nil {
		log.Error().Err(err).Msg("Get interim item counts failed")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, fmt.Errorf("error getting interim item count: %w", err)
	}

	// Save (clientID#dataType, mtime) into cache after writing into DB.
	for dataType, mtime := range typeMtimeMap {
		cache.SetTypeMtime(ctx, clientID, dataType, mtime)
	}

	err = db.UpdateClientItemCount(ctx, itemCounts, newNormalCount, newHistoryCount)
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

// handleClearServerDataRequest handles clearing user data from the datastore and cache
// and fills the response
func handleClearServerDataRequest(
	ctx context.Context,
	cache *cache.Cache,
	db datastore.Datastore,
	_ *sync_pb.ClearServerDataMessage,
	clientID string,
) (*sync_pb.SyncEnums_ErrorType, error) {
	errCode := sync_pb.SyncEnums_SUCCESS
	var err error

	err = db.DisableSyncChain(ctx, clientID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to disable sync chain")
		errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
		return &errCode, err
	}

	log.Info().Str("chainID", clientID).Msg("Clearing server data")

	syncEntities, err := db.ClearServerData(ctx, clientID)
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
		err = cache.Del(ctx, typeMtimeCacheKeys...)
		if err != nil {
			log.Error().Err(err).Msg("Failed to clear cache")
			errCode = sync_pb.SyncEnums_TRANSIENT_ERROR
			return &errCode, err
		}
	}

	log.Info().Str("chainID", clientID).Int("deletedCount", len(syncEntities)).Msg("Server data cleared")

	return &errCode, nil
}

func applyHandlerError(pbRsp *sync_pb.ClientToServerResponse, err error, wrapMsg string) error {
	if err == nil {
		return nil
	}
	if pbRsp.ErrorCode != nil {
		pbRsp.ErrorMessage = aws.String(err.Error())
		return nil
	}
	// In seldom error cases which are not temporary and will not go away
	// when clients retry, we will not use defined sync error in the proto
	// response, but use internal server error.
	return fmt.Errorf("%s: %w", wrapMsg, err)
}

// HandleClientToServerMessage handles the protobuf ClientToServerMessage and
// fills the protobuf ClientToServerResponse.
func HandleClientToServerMessage(
	ctx context.Context,
	cache *cache.Cache,
	pb *sync_pb.ClientToServerMessage,
	pbRsp *sync_pb.ClientToServerResponse,
	db datastore.Datastore,
	clientID string,
) error {
	// Create ClientToServerResponse and fill general fields for both GU and
	// Commit.
	pbRsp.StoreBirthday = aws.String(storeBirthday)
	pbRsp.ClientCommand = &sync_pb.ClientCommand{
		SetSyncPollInterval: aws.Int32(setSyncPollInterval),
		MaxCommitBatchSize:  aws.Int32(maxCommitBatchSize)}

	if pb.MessageContents == nil {
		return errors.New("nil pb.MessageContents received")
	}

	switch *pb.MessageContents {
	case sync_pb.ClientToServerMessage_GET_UPDATES:
		guRsp := &sync_pb.GetUpdatesResponse{}
		pbRsp.GetUpdates = guRsp
		var err error
		pbRsp.ErrorCode, err = handleGetUpdatesRequest(ctx, cache, pb.GetUpdates, guRsp, db, clientID)
		return applyHandlerError(pbRsp, err, "error handling GetUpdates request")
	case sync_pb.ClientToServerMessage_COMMIT:
		commitRsp := &sync_pb.CommitResponse{}
		pbRsp.Commit = commitRsp
		var err error
		pbRsp.ErrorCode, err = handleCommitRequest(context.TODO(), cache, pb.Commit, commitRsp, db, clientID)
		return applyHandlerError(pbRsp, err, "error handling Commit request")
	case sync_pb.ClientToServerMessage_CLEAR_SERVER_DATA:
		csdRsp := &sync_pb.ClearServerDataResponse{}
		pbRsp.ClearServerData = csdRsp
		var err error
		pbRsp.ErrorCode, err = handleClearServerDataRequest(
			context.Background(),
			cache,
			db,
			pb.ClearServerData,
			clientID,
		)
		return applyHandlerError(pbRsp, err, "error handling ClearServerData request")
	case sync_pb.ClientToServerMessage_DEPRECATED_3, sync_pb.ClientToServerMessage_DEPRECATED_4:
		fallthrough
	default:
		return errors.New("unsupported message type of ClientToServerMessage")
	}
}
