package command_test

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/suite"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

const (
	clientID     string = "client"
	bookmarkType int32  = 32904
	nigoriType   int32  = 47745
	cacheGUID    string = "cache_guid"
)

type CommandTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
	cache  *cache.Cache
}

type PBSyncAttrs struct {
	Name                   *string
	Version                *int64
	Deleted                *bool
	Folder                 *bool
	ServerDefinedUniqueTag *string
	Specifics              *sync_pb.EntitySpecifics
}

type PBSyncAttrsByName []*PBSyncAttrs

func (a PBSyncAttrsByName) Len() int           { return len(a) }
func (a PBSyncAttrsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PBSyncAttrsByName) Less(i, j int) bool { return *a[i].Name < *a[j].Name }

func NewPBSyncAttrs(name *string, version *int64, deleted *bool, folder *bool, serverTag *string, specifics *sync_pb.EntitySpecifics) *PBSyncAttrs {
	return &PBSyncAttrs{
		Name:                   name,
		Version:                version,
		Deleted:                deleted,
		Folder:                 folder,
		ServerDefinedUniqueTag: serverTag,
		Specifics:              specifics,
	}
}

func (suite *CommandTestSuite) SetupSuite() {
	datastore.Table = "client-entity-test-command"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")

	suite.cache = cache.NewCache(cache.NewRedisClient())
}

func (suite *CommandTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *CommandTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
	suite.Require().NoError(
		suite.cache.FlushAll(context.Background()), "Failed to clear cache")
}

func getNigoriSpecifics() *sync_pb.EntitySpecifics {
	nigoriEntitySpecifics := &sync_pb.EntitySpecifics_Nigori{
		Nigori: &sync_pb.NigoriSpecifics{},
	}
	return &sync_pb.EntitySpecifics{
		SpecificsVariant: nigoriEntitySpecifics,
	}
}

func getBookmarkSpecifics() *sync_pb.EntitySpecifics {
	bookmarkEntitySpecifics := &sync_pb.EntitySpecifics_Bookmark{
		Bookmark: &sync_pb.BookmarkSpecifics{},
	}
	return &sync_pb.EntitySpecifics{
		SpecificsVariant: bookmarkEntitySpecifics,
	}
}

func getCommitEntity(id string, version int64, deleted bool, specifics *sync_pb.EntitySpecifics) *sync_pb.SyncEntity {
	return &sync_pb.SyncEntity{
		IdString:  aws.String(id),
		Name:      aws.String(id),
		Version:   aws.Int64(version),
		Deleted:   aws.Bool(deleted),
		Folder:    aws.Bool(false),
		Specifics: specifics,
	}
}

func getClientCommand() *sync_pb.ClientCommand {
	return &sync_pb.ClientCommand{
		SetSyncPollInterval: aws.Int32(command.SetSyncPollInterval),
		MaxCommitBatchSize:  aws.Int32(command.MaxCommitBatchSize),
	}
}

func getClientToServerCommitMsg(entries []*sync_pb.SyncEntity) *sync_pb.ClientToServerMessage {
	commitMsg := &sync_pb.CommitMessage{
		Entries:   entries,
		CacheGuid: aws.String(cacheGUID),
	}
	contents := sync_pb.ClientToServerMessage_COMMIT
	return &sync_pb.ClientToServerMessage{
		MessageContents: &contents,
		Commit:          commitMsg,
	}
}

func getMarker(suite *CommandTestSuite, tokens []int64) []*sync_pb.DataTypeProgressMarker {
	types := []int32{nigoriType, bookmarkType} // hard-coded types used in tests.
	suite.Len(tokens, len(types))
	marker := []*sync_pb.DataTypeProgressMarker{}
	for i, token := range tokens {
		tokenBytes := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(tokenBytes, token)
		marker = append(marker, &sync_pb.DataTypeProgressMarker{
			DataTypeId: aws.Int32(types[i]), Token: tokenBytes})
	}
	return marker
}

func getClientToServerGUMsg(marker []*sync_pb.DataTypeProgressMarker,
	origin sync_pb.SyncEnums_GetUpdatesOrigin, fetchFolders bool,
	_ *int32) *sync_pb.ClientToServerMessage {
	guMsg := &sync_pb.GetUpdatesMessage{
		FetchFolders:       aws.Bool(fetchFolders),
		FromProgressMarker: marker,
		GetUpdatesOrigin:   &origin,
	}
	contents := sync_pb.ClientToServerMessage_GET_UPDATES
	return &sync_pb.ClientToServerMessage{
		MessageContents: &contents,
		GetUpdates:      guMsg,
	}
}

func getTokensFromNewMarker(suite *CommandTestSuite, newMarker []*sync_pb.DataTypeProgressMarker) (int64, int64) {
	nigoriToken, n := binary.Varint(newMarker[0].Token)
	suite.Positive(n)
	bookmarkToken, n := binary.Varint(newMarker[1].Token)
	suite.Positive(n)
	return nigoriToken, bookmarkToken
}

func assertCommonResponse(suite *CommandTestSuite, rsp *sync_pb.ClientToServerResponse, isCommit bool) {
	suite.Equal(sync_pb.SyncEnums_SUCCESS, *rsp.ErrorCode, "errorCode should match")
	suite.Equal(getClientCommand(), rsp.ClientCommand, "ClientCommand should match")
	suite.Equal(command.StoreBirthday, *rsp.StoreBirthday, "Birthday should match")
	if isCommit {
		suite.NotNil(rsp.Commit)
	} else {
		suite.NotNil(rsp.GetUpdates)
	}
}

func assertGetUpdatesResponse(suite *CommandTestSuite, rsp *sync_pb.GetUpdatesResponse,
	newMarker *[]*sync_pb.DataTypeProgressMarker, expectedPBSyncAttrs []*PBSyncAttrs,
	expectedChangesRemaining int64) {
	PBSyncAttrs := []*PBSyncAttrs{}
	for _, entity := range rsp.Entries {
		// Update tokens in the expected NewProgressMarker
		var tokenPtr *[]byte
		if strings.Contains(strings.ToLower(*entity.Name), "nigori") {
			tokenPtr = &(*newMarker)[0].Token
		} else { // bookmark type
			tokenPtr = &(*newMarker)[1].Token
		}
		token, n := binary.Varint(*tokenPtr)
		suite.Positive(n)
		if token < *entity.Mtime {
			binary.PutVarint(*tokenPtr, *entity.Mtime)
		}

		PBSyncAttrs = append(PBSyncAttrs,
			NewPBSyncAttrs(entity.Name, entity.Version, entity.Deleted,
				entity.Folder, entity.ServerDefinedUniqueTag, entity.Specifics))
	}

	sort.Sort(PBSyncAttrsByName(expectedPBSyncAttrs))
	sort.Sort(PBSyncAttrsByName(PBSyncAttrs))

	// Marshal to json to ignore protobuf internal fields when checking equality.
	s1, err := json.Marshal(expectedPBSyncAttrs)
	suite.Require().NoError(err, "json.Marshal should succeed")
	s2, err := json.Marshal(PBSyncAttrs)
	suite.Require().NoError(err, "json.Marshal should succeed")
	suite.Equal(s1, s2)

	suite.Equal(*newMarker, rsp.NewProgressMarker)
	suite.Equal(expectedChangesRemaining, *rsp.ChangesRemaining)
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_Basic() {
	// Prepare to commit 2 entries in 2 types.
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_nigori", 0, false, getNigoriSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	// Commit and check response.
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 2)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	serverIDs := []string{}
	commitVersions := []int64{}
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Equal(*entryRsp.Mtime, *entryRsp.Version)
		serverIDs = append(serverIDs, *entryRsp.IdString)
		commitVersions = append(commitVersions, *entryRsp.Version)
	}

	// GetUpdates with token 0 should get all of them.
	marker := getMarker(suite, []int64{0, 0})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)

	expectedPBSyncAttrs := []*PBSyncAttrs{
		NewPBSyncAttrs(entries[0].Name, &commitVersions[0], aws.Bool(false),
			aws.Bool(false), nil, getBookmarkSpecifics()),
		NewPBSyncAttrs(entries[1].Name, &commitVersions[1], aws.Bool(false),
			aws.Bool(false), nil, getNigoriSpecifics()),
	}
	newMarker := marker // Initialize expected NewProgressMarker with tokens = 0.
	assertGetUpdatesResponse(suite, rsp.GetUpdates, &newMarker, expectedPBSyncAttrs, 0)

	// Commit one new item, update one current item for each type.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity(serverIDs[0], commitVersions[0], true, getBookmarkSpecifics()),
		getCommitEntity(serverIDs[1], commitVersions[1], true, getNigoriSpecifics()),
		getCommitEntity("id3_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id4_nigori", 0, false, getNigoriSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)

	suite.Len(rsp.Commit.Entryresponse, 4)
	serverIDs = []string{}
	commitVersions = []int64{}
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Equal(*entryRsp.Mtime, *entryRsp.Version)
		serverIDs = append(serverIDs, *entryRsp.IdString)
		commitVersions = append(commitVersions, *entryRsp.Version)
	}

	// GetUpdates again with previous returned mtimes and check the result, it
	// should include update items and newly commit items.
	nigoriToken, bookmarkToken := getTokensFromNewMarker(suite, newMarker)
	marker = getMarker(suite, []int64{nigoriToken, bookmarkToken})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)

	expectedPBSyncAttrs = []*PBSyncAttrs{
		NewPBSyncAttrs(entries[0].Name, &commitVersions[0], aws.Bool(true),
			aws.Bool(false), nil, getBookmarkSpecifics()),
		NewPBSyncAttrs(entries[1].Name, &commitVersions[1], aws.Bool(true),
			aws.Bool(false), nil, getNigoriSpecifics()),
		NewPBSyncAttrs(entries[2].Name, &commitVersions[2], aws.Bool(false),
			aws.Bool(false), nil, getBookmarkSpecifics()),
		NewPBSyncAttrs(entries[3].Name, &commitVersions[3], aws.Bool(false),
			aws.Bool(false), nil, getNigoriSpecifics()),
	}
	newMarker = marker // Initialize expected NewProgressMarker with FromProgressMarker.
	assertGetUpdatesResponse(suite, rsp.GetUpdates, &newMarker, expectedPBSyncAttrs, 0)

	// Commit conflict items.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity(serverIDs[0], 1, false, getBookmarkSpecifics()),
		getCommitEntity(serverIDs[1], 1, false, getNigoriSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 2)
	commitConflict := sync_pb.CommitResponse_CONFLICT
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitConflict, *entryRsp.ResponseType)
	}

	// GetUpdates again with previous returned tokens should return 0 updates.
	nigoriToken, bookmarkToken = getTokensFromNewMarker(suite, newMarker)
	marker = getMarker(suite, []int64{nigoriToken, bookmarkToken})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)

	expectedPBSyncAttrs = []*PBSyncAttrs{}
	newMarker = marker // Initialize expected NewProgressMarker with FromProgressMarker.
	assertGetUpdatesResponse(suite, rsp.GetUpdates, &newMarker, expectedPBSyncAttrs, 0)
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_NewClient() {
	// Prepare input message for NEW_CLIENT get updates request.
	marker := getMarker(suite, []int64{0, 0})
	msg := getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_NEW_CLIENT, true, nil)
	rsp := &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)

	// Prepare expected sync entity attributes.
	expectedPBSyncAttrs := []*PBSyncAttrs{
		NewPBSyncAttrs(aws.String(command.NigoriName), aws.Int64(1),
			aws.Bool(false), aws.Bool(true), aws.String(command.NigoriTag),
			getNigoriSpecifics()),
		NewPBSyncAttrs(aws.String(command.BookmarksName), aws.Int64(1),
			aws.Bool(false), aws.Bool(true), aws.String(command.BookmarksTag),
			getBookmarkSpecifics()),
		NewPBSyncAttrs(aws.String(command.OtherBookmarksName), aws.Int64(1),
			aws.Bool(false), aws.Bool(true), aws.String(command.OtherBookmarksTag),
			getBookmarkSpecifics()),
		NewPBSyncAttrs(aws.String(command.SyncedBookmarksName), aws.Int64(1),
			aws.Bool(false), aws.Bool(true), aws.String(command.SyncedBookmarksTag),
			getBookmarkSpecifics()),
		NewPBSyncAttrs(aws.String(command.BookmarkBarName), aws.Int64(1),
			aws.Bool(false), aws.Bool(true), aws.String(command.BookmarkBarTag),
			getBookmarkSpecifics()),
	}
	newMarker := marker // Initialize expected NewProgressMarker with tokens = 0.
	assertGetUpdatesResponse(suite, rsp.GetUpdates, &newMarker, expectedPBSyncAttrs, 0)

	// Check dummy encryption keys only for NEW_CLIENT case.
	expectedEncryptionKeys := make([][]byte, 1)
	expectedEncryptionKeys[0] = []byte("1234")
	suite.Equal(expectedEncryptionKeys, rsp.GetUpdates.EncryptionKeys)
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_GUBatchSize() {
	// Commit a few items for testing.
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_nigori", 0, false, getNigoriSpecifics()),
		getCommitEntity("id3_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id4_nigori", 0, false, getNigoriSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	// Commit and check response.
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 4)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Equal(*entryRsp.Mtime, *entryRsp.Version)
	}
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_QuotaLimit() {
	defaultMaxClientObjectQuota := *command.MaxClientObjectQuota
	*command.MaxClientObjectQuota = 4

	// Commit 2 items without exceed quota.
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 2)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	serverIDs := []string{}
	commitVersions := []int64{}
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Equal(*entryRsp.Mtime, *entryRsp.Version)
		serverIDs = append(serverIDs, *entryRsp.IdString)
		commitVersions = append(commitVersions, *entryRsp.Version)
	}

	// Commit 4 items to exceed quota by a half, 2 should return OVER_QUOTA.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity("id3_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id4_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id5_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id6_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 4)
	overQuota := sync_pb.CommitResponse_OVER_QUOTA
	expectedEntryRsp := []sync_pb.CommitResponse_ResponseType{commitSuccess, commitSuccess, overQuota, overQuota}
	expectedVersion := []*int64{rsp.Commit.Entryresponse[0].Mtime, rsp.Commit.Entryresponse[1].Mtime, nil, nil}
	for i, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(expectedEntryRsp[i], *entryRsp.ResponseType)
		suite.Equal(expectedVersion[i], entryRsp.Version)
	}

	// Commit 2 items again when quota is already exceed should get two OVER_QUOTA
	entries = []*sync_pb.SyncEntity{
		getCommitEntity("id7_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id8_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 2)
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(overQuota, *entryRsp.ResponseType)
	}

	// Commit updates to delete two previous inserted items.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity(serverIDs[0], commitVersions[0], true, getBookmarkSpecifics()),
		getCommitEntity(serverIDs[1], commitVersions[1], true, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 2)
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Equal(*entryRsp.Mtime, *entryRsp.Version)
	}

	// Commit 4 items should have two success and two OVER_QUOTA.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity("id7_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id8_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id9_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id10_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 4)
	expectedVersion = []*int64{rsp.Commit.Entryresponse[0].Mtime, rsp.Commit.Entryresponse[1].Mtime, nil, nil}
	for i, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(expectedEntryRsp[i], *entryRsp.ResponseType)
		if *entryRsp.ResponseType == commitSuccess {
			suite.Equal(*expectedVersion[i], *entryRsp.Version)
		} else {
			suite.Equal(expectedVersion[i], entryRsp.Version)
		}
	}

	*command.MaxClientObjectQuota = defaultMaxClientObjectQuota
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_ReplaceParentIDToServerGeneratedID() {
	child0 := getCommitEntity("id_child0", 0, false, getBookmarkSpecifics())
	msg := getClientToServerCommitMsg([]*sync_pb.SyncEntity{child0})
	rsp := &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 1)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
	}

	// Commit parents with its child bookmarks in one commit request.
	parent1 := getCommitEntity("id_parent", 0, false, getBookmarkSpecifics())
	parent1.Folder = aws.Bool(true)
	child1 := getCommitEntity("id_child", 0, false, getBookmarkSpecifics())
	child1.ParentIdString = aws.String("id_parent")
	parent2 := getCommitEntity("id_parent2", 0, false, getBookmarkSpecifics())
	parent2.Folder = aws.Bool(true)
	child2 := getCommitEntity("id_child2", 0, false, getBookmarkSpecifics())
	child2.ParentIdString = aws.String("id_parent")
	child3 := getCommitEntity("id_child3", 0, false, getBookmarkSpecifics())
	child3.ParentIdString = aws.String("id_parent2")

	updateChild0 := getCommitEntity(*rsp.Commit.Entryresponse[0].IdString, *rsp.Commit.Entryresponse[0].Version, false, getBookmarkSpecifics())
	updateChild0.ParentIdString = aws.String("id_parent")

	entries := []*sync_pb.SyncEntity{parent1, child1, parent2, child2, child3, updateChild0}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 6)
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
	}

	// Get updates to check if child's parent ID is replaced with the server
	// generated ID of its parent.
	marker := getMarker(suite, []int64{0, 0})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, true, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Require().Len(rsp.GetUpdates.Entries, 6)
	for i := range rsp.GetUpdates.Entries {
		if i != len(rsp.GetUpdates.Entries)-1 {
			suite.Equal(rsp.GetUpdates.Entries[i].OriginatorClientItemId, entries[i].IdString)
		} else {
			suite.Equal(rsp.GetUpdates.Entries[i].OriginatorClientItemId, child0.IdString)
		}
		suite.NotNil(rsp.GetUpdates.Entries[i].IdString)
	}

	suite.Equal(rsp.GetUpdates.Entries[1].ParentIdString, rsp.GetUpdates.Entries[0].IdString)
	suite.Equal(rsp.GetUpdates.Entries[3].ParentIdString, rsp.GetUpdates.Entries[0].IdString)
	suite.Equal(rsp.GetUpdates.Entries[4].ParentIdString, rsp.GetUpdates.Entries[2].IdString)
	suite.Equal(rsp.GetUpdates.Entries[5].ParentIdString, rsp.GetUpdates.Entries[0].IdString)
}

func assertTypeMtimeCacheValue(suite *CommandTestSuite, key string, mtime int64, errMsg string) {
	val, err := suite.cache.Get(context.Background(), key, false)
	suite.Require().NoError(err, "cache.Get should succeed")
	suite.Equal(val, strconv.FormatInt(mtime, 10), errMsg)
}

func insertSyncEntitiesWithoutUpdateCache(
	suite *CommandTestSuite, entries []*sync_pb.SyncEntity, clientID string) (ret []*datastore.SyncEntity) {
	for _, entry := range entries {
		dbEntry, err := datastore.CreateDBSyncEntity(entry, nil, clientID)
		suite.Require().NoError(err, "Create db entity from pb entity should succeed")
		_, err = suite.dynamo.InsertSyncEntity(dbEntry)
		suite.Require().NoError(err, "Insert sync entity should succeed")
		val, err := suite.cache.Get(context.Background(),
			clientID+"#"+strconv.Itoa(*dbEntry.DataType), false)
		suite.Require().NoError(err, "Get from cache should succeed")
		suite.Require().NotEqual(val, strconv.FormatInt(*dbEntry.Mtime, 10),
			"Cache should not be updated")
		ret = append(ret, dbEntry)
	}
	return
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_TypeMtimeCache_Basic() {
	// Commit two entries of type1, one entry of type2 into dynamoDB and get the
	// mtime from response.
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id3_nigori", 0, false, getNigoriSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 3)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	var latestBookmarkMtime int64
	var latestNigoriMtime int64
	for i, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		if i < 2 {
			latestBookmarkMtime = *entryRsp.Mtime
		}
		if i == 2 {
			latestNigoriMtime = *entryRsp.Mtime
		}
	}

	// Latest mtime of each type in the commit should be stored in the cache.
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime,
		"Successful commit should write the latest mtime into cache")
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(nigoriType)),
		latestNigoriMtime,
		"Successful commit should write the latest mtime into cache")

	// Insert an entry into DB manually to make sure there are updates in DB
	// after lastestBookmark time to check if we do short circuit in later GU.
	dbEntries := insertSyncEntitiesWithoutUpdateCache(suite,
		[]*sync_pb.SyncEntity{
			getCommitEntity("id4_bookmark", 0, false, getBookmarkSpecifics()),
			getCommitEntity("id5_nigori", 0, false, getNigoriSpecifics()),
		},
		clientID)

	// GU request with the same or newer token should be short circuited, so
	// should return no updates.
	marker := getMarker(suite, []int64{latestNigoriMtime, latestBookmarkMtime + 1})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_PERIODIC, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Empty(rsp.GetUpdates.Entries)
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime, "cache is not updated when short circuited")
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(nigoriType)),
		latestNigoriMtime, "cache is not updated when short circuited")

	// Manually update cache for our DB insert.
	latestBookmarkMtime = *dbEntries[0].Mtime
	latestNigoriMtime = *dbEntries[1].Mtime
	suite.cache.SetTypeMtime(context.Background(), clientID, int(bookmarkType), latestBookmarkMtime)
	suite.cache.SetTypeMtime(context.Background(), clientID, int(nigoriType), latestNigoriMtime)

	// Commit another entry and check if cache is updated.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity("id6_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 1)
	entryRsp := rsp.Commit.Entryresponse[0]
	suite.Equal(commitSuccess, *entryRsp.ResponseType)

	latestBookmarkMtime = *entryRsp.Mtime
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime, "Successful commit should update the cache")

	// Send GU with an old token will get updates immediately.
	// Check the cache value again, should be the same as the latest mtime in rsp.
	marker = getMarker(suite, []int64{latestNigoriMtime - 1, latestBookmarkMtime - 1})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_PERIODIC, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Len(rsp.GetUpdates.Entries, 2)
	suite.Equal(latestNigoriMtime, *rsp.GetUpdates.Entries[0].Mtime)
	suite.Equal(latestBookmarkMtime, *rsp.GetUpdates.Entries[1].Mtime)
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime, "Cached token should be equal to latest mtime")
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(nigoriType)),
		latestNigoriMtime, "Cached token should be equal to latest mtime")
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_TypeMtimeCache_SkipCacheForNonPollReq() {
	// Commit one entity and check cache value is set properly.
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 1)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	suite.Equal(commitSuccess, *rsp.Commit.Entryresponse[0].ResponseType)
	latestBookmarkMtime := *rsp.Commit.Entryresponse[0].Mtime
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime,
		"Commit should write the latest mtime into cache")

	// Make sure non-poll request should not be short circuited due to cache.
	// Insert an entry into DB manually without touching the cache to make sure
	// there are updates in DB after lastestBookmark so we will have updates if
	// we go query the DB using the previous token.
	dbEntries := insertSyncEntitiesWithoutUpdateCache(suite,
		[]*sync_pb.SyncEntity{
			getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
		},
		clientID)

	// Check that we will receive the manually inserted item from DB immediately.
	marker := getMarker(suite, []int64{0, latestBookmarkMtime})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, true, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Require().Len(rsp.GetUpdates.Entries, 1)
	suite.Require().Equal(dbEntries[0].Mtime, rsp.GetUpdates.Entries[0].Mtime)
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		*dbEntries[0].Mtime, "Successful commit should update the cache")
}

func (suite *CommandTestSuite) TestHandleClientToServerMessage_TypeMtimeCache_ChangesRemaining() {
	// Commit two entries and check cache value is set properly.
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Len(rsp.Commit.Entryresponse, 2)
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	var latestBookmarkMtime int64
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Equal(commitSuccess, *entryRsp.ResponseType)
		suite.NotEqual(latestBookmarkMtime, *entryRsp.Mtime)
		latestBookmarkMtime = *entryRsp.Mtime
	}
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime,
		"Commit should write the latest mtime into cache")

	// Send a GU with batch size set to 1, changesRemaining in rsp should be 1
	// and cache should not be updated.
	marker := getMarker(suite, []int64{0, 0})
	clientBatch := int32(2)
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_PERIODIC, true, &clientBatch)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Require().Len(rsp.GetUpdates.Entries, 2)
	suite.Require().Equal(int64(0), *rsp.GetUpdates.ChangesRemaining)
	mtime := *rsp.GetUpdates.Entries[0].Mtime
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime,
		"cache should not be updated when changes remaining = 1")

	// Send a second GU with changesRemaining in rsp = 0 and check cache is now
	// updated.
	marker = getMarker(suite, []int64{0, mtime})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_PERIODIC, true, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamo, clientID),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Require().Len(rsp.GetUpdates.Entries, 1)
	suite.Require().Equal(int64(0), *rsp.GetUpdates.ChangesRemaining)
	assertTypeMtimeCacheValue(suite, clientID+"#"+strconv.Itoa(int(bookmarkType)),
		latestBookmarkMtime,
		"cache should be updated when changes remaining = 0")
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
