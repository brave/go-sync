package command_test

import (
	"encoding/binary"
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/stretchr/testify/suite"
)

const (
	bookmarkType int32  = 32904
	nigoriType   int32  = 47745
	cacheGUID    string = "cache_guid"
)

type CommandTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
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
}

func (suite *CommandTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *CommandTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
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
		SetSyncPollInterval:        aws.Int32(command.SetSyncPollInterval),
		MaxCommitBatchSize:         aws.Int32(command.MaxCommitBatchSize),
		SessionsCommitDelaySeconds: aws.Int32(command.SessionsCommitDelaySeconds),
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
	suite.Assert().Equal(len(types), len(tokens))
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
	batchSize *int32) *sync_pb.ClientToServerMessage {
	guMsg := &sync_pb.GetUpdatesMessage{
		FetchFolders:       aws.Bool(fetchFolders),
		FromProgressMarker: marker,
		GetUpdatesOrigin:   &origin,
		BatchSize:          batchSize,
	}
	contents := sync_pb.ClientToServerMessage_GET_UPDATES
	return &sync_pb.ClientToServerMessage{
		MessageContents: &contents,
		GetUpdates:      guMsg,
	}
}

func getTokensFromNewMarker(suite *CommandTestSuite, newMarker []*sync_pb.DataTypeProgressMarker) (int64, int64) {
	nigoriToken, n := binary.Varint(newMarker[0].Token)
	suite.Assert().Greater(n, 0)
	bookmarkToken, n := binary.Varint(newMarker[1].Token)
	suite.Assert().Greater(n, 0)
	return nigoriToken, bookmarkToken
}

func assertCommonResponse(suite *CommandTestSuite, rsp *sync_pb.ClientToServerResponse, isCommit bool) {
	suite.Assert().Equal(sync_pb.SyncEnums_SUCCESS, *rsp.ErrorCode, "errorCode should match")
	suite.Assert().Equal(getClientCommand(), rsp.ClientCommand, "ClientCommand should match")
	suite.Assert().Equal(command.StoreBirthday, *rsp.StoreBirthday, "Birthday should match")
	if isCommit {
		suite.Assert().NotNil(rsp.Commit)
	} else {
		suite.Assert().NotNil(rsp.GetUpdates)
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
		suite.Assert().Greater(n, 0)
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
	suite.Assert().Equal(s1, s2)

	suite.Assert().Equal(*newMarker, rsp.NewProgressMarker)
	suite.Assert().Equal(expectedChangesRemaining, *rsp.ChangesRemaining)
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(2, len(rsp.Commit.Entryresponse))
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	serverIDs := []string{}
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Assert().Equal(int64(1), *entryRsp.Version)
		serverIDs = append(serverIDs, *entryRsp.IdString)
	}

	// GetUpdates with token 0 should get all of them.
	marker := getMarker(suite, []int64{0, 0})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)

	expectedPBSyncAttrs := []*PBSyncAttrs{
		NewPBSyncAttrs(entries[0].Name, aws.Int64(1), aws.Bool(false),
			aws.Bool(false), nil, getBookmarkSpecifics()),
		NewPBSyncAttrs(entries[1].Name, aws.Int64(1), aws.Bool(false),
			aws.Bool(false), nil, getNigoriSpecifics()),
	}
	newMarker := marker // Initialize expected NewProgressMarker with tokens = 0.
	assertGetUpdatesResponse(suite, rsp.GetUpdates, &newMarker, expectedPBSyncAttrs, 0)

	// Commit one new item, update one current item for each type.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity(serverIDs[0], 1, true, getBookmarkSpecifics()),
		getCommitEntity(serverIDs[1], 1, true, getNigoriSpecifics()),
		getCommitEntity("id3_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id4_nigori", 0, false, getNigoriSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)

	suite.Assert().Equal(4, len(rsp.Commit.Entryresponse))
	serverIDs = []string{}
	expectedVersion := []int64{2, 2, 1, 1}
	for i, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Assert().Equal(expectedVersion[i], *entryRsp.Version)
		serverIDs = append(serverIDs, *entryRsp.IdString)
	}

	// GetUpdates again with previous returned mtimes and check the result, it
	// should include update items and newly commit items.
	nigoriToken, bookmarkToken := getTokensFromNewMarker(suite, newMarker)
	marker = getMarker(suite, []int64{nigoriToken, bookmarkToken})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)

	expectedPBSyncAttrs = []*PBSyncAttrs{
		NewPBSyncAttrs(entries[0].Name, aws.Int64(2), aws.Bool(true),
			aws.Bool(false), nil, getBookmarkSpecifics()),
		NewPBSyncAttrs(entries[1].Name, aws.Int64(2), aws.Bool(true),
			aws.Bool(false), nil, getNigoriSpecifics()),
		NewPBSyncAttrs(entries[2].Name, aws.Int64(1), aws.Bool(false),
			aws.Bool(false), nil, getBookmarkSpecifics()),
		NewPBSyncAttrs(entries[3].Name, aws.Int64(1), aws.Bool(false),
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(2, len(rsp.Commit.Entryresponse))
	commitConflict := sync_pb.CommitResponse_CONFLICT
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(commitConflict, *entryRsp.ResponseType)
	}

	// GetUpdates again with previous returned tokens should return 0 updates.
	nigoriToken, bookmarkToken = getTokensFromNewMarker(suite, newMarker)
	marker = getMarker(suite, []int64{nigoriToken, bookmarkToken})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
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
	suite.Assert().Equal(expectedEncryptionKeys, rsp.GetUpdates.EncryptionKeys)
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(4, len(rsp.Commit.Entryresponse))
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Assert().Equal(int64(1), *entryRsp.Version)
	}

	// Test maxGUBatchSize from client side should be respected when smaller than
	// the server one.
	clientBatchSize := 3
	marker := getMarker(suite, []int64{0, 0})
	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, aws.Int32(int32(clientBatchSize)))
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Assert().Equal(clientBatchSize, len(rsp.GetUpdates.Entries))
	suite.Assert().Equal(
		int64(len(entries)-clientBatchSize), *rsp.GetUpdates.ChangesRemaining)
	// nigori1, nigori2, bookmark1
	expectedName := []*string{entries[1].Name, entries[3].Name, entries[0].Name}
	for i, entry := range rsp.GetUpdates.Entries {
		suite.Assert().Equal(expectedName[i], entry.Name)
	}
	expectedNewMarker := getMarker(suite,
		[]int64{*rsp.GetUpdates.Entries[1].Mtime, *rsp.GetUpdates.Entries[2].Mtime})
	suite.Assert().Equal(expectedNewMarker, rsp.GetUpdates.NewProgressMarker)

	// Test maxGUBatchSize from server side should be respected when smaller than
	// the client one.
	defaultServerGUBatchSize := *command.MaxGUBatchSize
	*command.MaxGUBatchSize = 2
	rsp = &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, false)
	suite.Assert().Equal(int(*command.MaxGUBatchSize), len(rsp.GetUpdates.Entries))
	suite.Assert().Equal(int64(1), *rsp.GetUpdates.ChangesRemaining)
	// nigori1, nigori2
	expectedName = []*string{entries[1].Name, entries[3].Name}
	for i, entry := range rsp.GetUpdates.Entries {
		suite.Assert().Equal(expectedName[i], entry.Name)
	}
	expectedNewMarker = getMarker(suite, []int64{*rsp.GetUpdates.Entries[1].Mtime, 0})
	suite.Assert().Equal(expectedNewMarker, rsp.GetUpdates.NewProgressMarker)
	*command.MaxGUBatchSize = defaultServerGUBatchSize
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(2, len(rsp.Commit.Entryresponse))
	commitSuccess := sync_pb.CommitResponse_SUCCESS
	serverIDs := []string{}
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Assert().Equal(int64(1), *entryRsp.Version)
		serverIDs = append(serverIDs, *entryRsp.IdString)
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(4, len(rsp.Commit.Entryresponse))
	overQuota := sync_pb.CommitResponse_OVER_QUOTA
	expectedEntryRsp := []sync_pb.CommitResponse_ResponseType{commitSuccess, commitSuccess, overQuota, overQuota}
	expectedVersion := []*int64{aws.Int64(1), aws.Int64(1), nil, nil}
	for i, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(expectedEntryRsp[i], *entryRsp.ResponseType)
		suite.Assert().Equal(expectedVersion[i], entryRsp.Version)
	}

	// Commit 2 items again when quota is already exceed should get two OVER_QUOTA
	entries = []*sync_pb.SyncEntity{
		getCommitEntity("id7_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id8_bookmark", 0, false, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(2, len(rsp.Commit.Entryresponse))
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(overQuota, *entryRsp.ResponseType)
	}

	// Commit updates to delete two previous inserted items.
	entries = []*sync_pb.SyncEntity{
		getCommitEntity(serverIDs[0], 1, true, getBookmarkSpecifics()),
		getCommitEntity(serverIDs[1], 1, true, getBookmarkSpecifics()),
	}
	msg = getClientToServerCommitMsg(entries)
	rsp = &sync_pb.ClientToServerResponse{}

	suite.Require().NoError(
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(2, len(rsp.Commit.Entryresponse))
	for _, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(commitSuccess, *entryRsp.ResponseType)
		suite.Assert().Equal(int64(2), *entryRsp.Version)
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
		command.HandleClientToServerMessage(msg, rsp, suite.dynamo, "client"),
		"HandleClientToServerMessage should succeed")
	assertCommonResponse(suite, rsp, true)
	suite.Assert().Equal(4, len(rsp.Commit.Entryresponse))
	for i, entryRsp := range rsp.Commit.Entryresponse {
		suite.Assert().Equal(expectedEntryRsp[i], *entryRsp.ResponseType)
		suite.Assert().Equal(expectedVersion[i], entryRsp.Version)
	}

	*command.MaxClientObjectQuota = defaultMaxClientObjectQuota
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
