package command_test

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

const (
	getUpdatesNigoriType     int32  = 47745
	getUpdatesBookmarkType   int32  = 32904
	getUpdatesDeviceInfoType int    = 154522
	getUpdatesNigoriTag      string = "google_chrome_nigori"
	getUpdatesDefaultMaxSize int    = 500
)

type GetUpdatesTestSuite struct {
	suite.Suite

	cache *cache.Cache
}

func (suite *GetUpdatesTestSuite) SetupSuite() {
	suite.cache = cache.NewCache(cache.NewRedisClient())
}

func (suite *GetUpdatesTestSuite) TearDownTest() {
	suite.Require().NoError(
		suite.cache.FlushAll(context.Background()), "Failed to clear cache")
}

func tokenBytes(value int64) []byte {
	token := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(token, value)
	return token
}

func getUpdatesRequest(
	marker []*sync_pb.DataTypeProgressMarker,
	origin *sync_pb.SyncEnums_GetUpdatesOrigin,
) *sync_pb.ClientToServerMessage {
	contents := sync_pb.ClientToServerMessage_GET_UPDATES
	return &sync_pb.ClientToServerMessage{
		MessageContents: &contents,
		GetUpdates: &sync_pb.GetUpdatesMessage{
			FromProgressMarker: marker,
			GetUpdatesOrigin:   origin,
		},
	}
}

func (suite *GetUpdatesTestSuite) handleGetUpdates(
	msg *sync_pb.ClientToServerMessage,
	db datastore.Datastore,
) *sync_pb.ClientToServerResponse {
	rsp := &sync_pb.ClientToServerResponse{}
	err := command.HandleClientToServerMessage(
		context.Background(),
		suite.cache,
		msg,
		rsp,
		db,
		"client",
	)
	suite.Require().NoError(err)
	return rsp
}

func (suite *GetUpdatesTestSuite) TestHandleGetUpdatesRequest_NilProgressMarkers() {
	rsp := suite.handleGetUpdates(
		getUpdatesRequest(nil, nil),
		new(datastoretest.MockDatastore),
	)
	suite.Require().NotNil(rsp.ErrorCode)
	suite.Equal(sync_pb.SyncEnums_SUCCESS, *rsp.ErrorCode)
	suite.Require().NotNil(rsp.GetUpdates)
	suite.Equal(int64(0), *rsp.GetUpdates.ChangesRemaining)
}

func (suite *GetUpdatesTestSuite) TestHandleGetUpdatesRequest_InvalidToken() {
	rsp := &sync_pb.ClientToServerResponse{}
	err := command.HandleClientToServerMessage(
		context.Background(),
		suite.cache,
		getUpdatesRequest([]*sync_pb.DataTypeProgressMarker{{
			DataTypeId: aws.Int32(getUpdatesNigoriType),
			Token:      []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80},
		}}, nil),
		rsp,
		new(datastoretest.MockDatastore),
		"client",
	)
	suite.Require().Error(err)
	suite.Contains(err.Error(), "failed at decoding token value")
}

func (suite *GetUpdatesTestSuite) TestHandleGetUpdatesRequest_CreatePBSyncEntityError() {
	mtime := int64(1)
	db := new(datastoretest.MockDatastore)
	db.On("GetUpdatesForType", mock.Anything, int(getUpdatesNigoriType), int64(0), true, "client", getUpdatesDefaultMaxSize).
		Return(false, []datastore.SyncEntity{{
			ID:        "nigori",
			Version:   aws.Int64(1),
			Mtime:     &mtime,
			Specifics: []byte("not-valid-protobuf"),
		}}, nil)

	rsp := suite.handleGetUpdates(
		getUpdatesRequest([]*sync_pb.DataTypeProgressMarker{{
			DataTypeId: aws.Int32(getUpdatesNigoriType),
			Token:      tokenBytes(0),
		}}, nil),
		db,
	)
	suite.Require().NotNil(rsp.ErrorCode)
	suite.Equal(sync_pb.SyncEnums_TRANSIENT_ERROR, *rsp.ErrorCode)
	suite.Require().NotNil(rsp.ErrorMessage)
	suite.Contains(*rsp.ErrorMessage, "error creating protobuf sync entity")
	db.AssertExpectations(suite.T())
}

func (suite *GetUpdatesTestSuite) TestHandleGetUpdatesRequest_GetUpdatesForTypeError() {
	db := new(datastoretest.MockDatastore)
	db.On("GetUpdatesForType", mock.Anything, int(getUpdatesNigoriType), int64(0), true, "client", getUpdatesDefaultMaxSize).
		Return(false, []datastore.SyncEntity{}, errors.New("db down"))

	rsp := suite.handleGetUpdates(
		getUpdatesRequest([]*sync_pb.DataTypeProgressMarker{{
			DataTypeId: aws.Int32(getUpdatesNigoriType),
			Token:      tokenBytes(0),
		}}, nil),
		db,
	)
	suite.Require().NotNil(rsp.ErrorCode)
	suite.Equal(sync_pb.SyncEnums_TRANSIENT_ERROR, *rsp.ErrorCode)
	db.AssertExpectations(suite.T())
}

func (suite *GetUpdatesTestSuite) TestHandleGetUpdatesRequest_NigoriRootNotReady() {
	db := new(datastoretest.MockDatastore)
	db.On("GetUpdatesForType", mock.Anything, getUpdatesDeviceInfoType, int64(0), false, "client", getUpdatesDefaultMaxSize).
		Return(false, []datastore.SyncEntity{}, nil)
	db.On("HasServerDefinedUniqueTag", mock.Anything, "client", getUpdatesNigoriTag).
		Return(true, nil)
	db.On("GetUpdatesForType", mock.Anything, int(getUpdatesNigoriType), int64(0), true, "client", getUpdatesDefaultMaxSize).
		Return(false, []datastore.SyncEntity{}, nil)

	origin := sync_pb.SyncEnums_NEW_CLIENT
	rsp := suite.handleGetUpdates(
		getUpdatesRequest([]*sync_pb.DataTypeProgressMarker{{
			DataTypeId: aws.Int32(getUpdatesNigoriType),
			Token:      tokenBytes(0),
		}}, &origin),
		db,
	)
	suite.Require().NotNil(rsp.ErrorCode)
	suite.Equal(sync_pb.SyncEnums_TRANSIENT_ERROR, *rsp.ErrorCode)
	suite.Require().NotNil(rsp.ErrorMessage)
	suite.Contains(*rsp.ErrorMessage, "nigori root folder entity is not ready")
}

func (suite *GetUpdatesTestSuite) TestHandleGetUpdatesRequest_SkipTypeWhenBatchFull() {
	mtime := int64(7)
	entities := make([]datastore.SyncEntity, getUpdatesDefaultMaxSize)
	for i := range entities {
		itemMtime := mtime + int64(i)
		entities[i] = datastore.SyncEntity{
			ID:      "nigori",
			Version: aws.Int64(1),
			Mtime:   &itemMtime,
			Ctime:   aws.Int64(1),
			Deleted: aws.Bool(false),
			Folder:  aws.Bool(true),
		}
	}
	db := new(datastoretest.MockDatastore)
	db.On("GetUpdatesForType", mock.Anything, int(getUpdatesNigoriType), int64(0), true, "client", getUpdatesDefaultMaxSize).
		Return(false, entities, nil)

	bookmarkToken := tokenBytes(3)
	rsp := suite.handleGetUpdates(
		getUpdatesRequest([]*sync_pb.DataTypeProgressMarker{
			{DataTypeId: aws.Int32(getUpdatesNigoriType), Token: tokenBytes(0)},
			{DataTypeId: aws.Int32(getUpdatesBookmarkType), Token: bookmarkToken},
		}, nil),
		db,
	)
	suite.Require().NotNil(rsp.ErrorCode)
	suite.Equal(sync_pb.SyncEnums_SUCCESS, *rsp.ErrorCode)
	suite.Require().NotNil(rsp.GetUpdates)
	suite.Require().Len(rsp.GetUpdates.Entries, getUpdatesDefaultMaxSize)
	suite.Require().Len(rsp.GetUpdates.NewProgressMarker, 2)
	suite.Equal(bookmarkToken, rsp.GetUpdates.NewProgressMarker[1].Token)
	db.AssertExpectations(suite.T())
}

func TestGetUpdatesTestSuite(t *testing.T) {
	suite.Run(t, new(GetUpdatesTestSuite))
}
