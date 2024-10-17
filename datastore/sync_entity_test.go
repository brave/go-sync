package datastore_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
)

type SyncEntityTestSuite struct {
	suite.Suite
}

func (suite *SyncEntityTestSuite) TestNewServerClientUniqueTagItem() {
	expectedServerTag := datastore.ServerClientUniqueTagItem{
		ClientID: "id",
		ID:       "Server#serverTag",
	}
	expectedClientTag := datastore.ServerClientUniqueTagItem{
		ClientID: "id",
		ID:       "Client#clientTag",
	}
	actualClientTag := *datastore.NewServerClientUniqueTagItem("id", "clientTag", false)
	actualServerTag := *datastore.NewServerClientUniqueTagItem("id", "serverTag", true)

	// We can't know the exact value for Mtime & Ctime.  Make sure they're set,
	// set zero value for subsequent tests
	suite.Assert().NotNil(actualClientTag.Mtime)
	suite.Assert().NotNil(actualClientTag.Ctime)
	suite.Assert().NotNil(actualServerTag.Mtime)
	suite.Assert().NotNil(actualServerTag.Ctime)

	actualClientTag.Mtime = nil
	actualClientTag.Ctime = nil
	actualServerTag.Mtime = nil
	actualServerTag.Ctime = nil

	suite.Assert().Equal(expectedServerTag, actualServerTag)
	suite.Assert().Equal(expectedClientTag, actualClientTag)
}

func (suite *SyncEntityTestSuite) TestCreateDBSyncEntity() {
	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	nigoriEntitySpecific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: nigoriEntitySpecific}
	specificsBytes, err := proto.Marshal(specifics)
	suite.Require().NoError(err, "Marshal specifics should succeed")

	uniquePosition := &sync_pb.UniquePosition{
		CustomCompressedV1: []byte{1, 2},
	}
	uniquePositionBytes, err := proto.Marshal(uniquePosition)
	suite.Require().NoError(err, "Marshal unique position should succeed")

	guid := aws.String("guid")
	pbEntity := sync_pb.SyncEntity{
		IdString:               aws.String("client_item_id"),
		ParentIdString:         aws.String("parent_id"),
		Version:                aws.Int64(0),
		Name:                   aws.String("name"),
		NonUniqueName:          aws.String("non_unique_name"),
		ServerDefinedUniqueTag: aws.String("server_tag"),
		ClientTagHash:          aws.String("client_tag"),
		Deleted:                aws.Bool(false),
		Folder:                 aws.Bool(false),
		Specifics:              specifics,
		UniquePosition:         uniquePosition,
	}
	var expectedChainID int64 = 1
	expectedDBEntity := datastore.SyncEntity{
		ClientID:               "client1",
		ChainID:                &expectedChainID,
		ParentID:               pbEntity.ParentIdString,
		Version:                pbEntity.Version,
		Name:                   pbEntity.Name,
		NonUniqueName:          pbEntity.NonUniqueName,
		ServerDefinedUniqueTag: pbEntity.ServerDefinedUniqueTag,
		ClientDefinedUniqueTag: pbEntity.ClientTagHash,
		Deleted:                pbEntity.Deleted,
		Folder:                 pbEntity.Folder,
		Specifics:              specificsBytes,
		UniquePosition:         uniquePositionBytes,
		DataType:               aws.Int(47745), // nigori type ID
		OriginatorCacheGUID:    guid,
		OriginatorClientItemID: pbEntity.IdString,
		ExpirationTime:         nil,
	}

	dbEntity, err := datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")

	// Check ID is replaced with a server-generated ID.
	suite.Assert().NotEqual(
		dbEntity.ID, *pbEntity.IdString,
		"ID should be a server-generated ID and not equal to the passed IdString")
	_, err = uuid.Parse(dbEntity.ID)
	suite.Assert().NoError(err, "dbEntity.ID should be a valid UUID")

	expectedDBEntity.ID = dbEntity.ID

	// Check Mtime and Ctime should be provided by the server if client does not
	// provide it.
	suite.Assert().NotNil(
		dbEntity.Ctime, "Mtime should not be nil if client did not pass one")
	suite.Assert().NotNil(
		dbEntity.Mtime, "Mtime should not be nil if client did not pass one")
	suite.Assert().Equal(
		*dbEntity.Mtime, *dbEntity.Ctime,
		"Server should generate the same value for mtime and ctime when they're not provided by the client")
	expectedDBEntity.Ctime = dbEntity.Ctime
	expectedDBEntity.Mtime = dbEntity.Mtime
	expectedDBEntity.DataTypeMtime = aws.String("47745#" + strconv.FormatInt(*dbEntity.Mtime, 10))
	suite.Assert().Equal(dbEntity, &expectedDBEntity)
	suite.Assert().Nil(dbEntity.ExpirationTime)

	pbEntity.Deleted = nil
	pbEntity.Folder = nil
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().False(*dbEntity.Deleted, "Default value should be set for Deleted for new entities")
	suite.Assert().False(*dbEntity.Folder, "Default value should be set for Deleted for new entities")
	suite.Assert().Nil(dbEntity.ExpirationTime)

	// Check the case when Ctime and Mtime are provided by the client.
	pbEntity.Ctime = aws.Int64(12345678)
	pbEntity.Mtime = aws.Int64(12345678)
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().Equal(*dbEntity.Ctime, *pbEntity.Ctime, "Client's Ctime should be respected")
	suite.Assert().NotEqual(*dbEntity.Mtime, *pbEntity.Mtime, "Client's Mtime should be replaced")
	suite.Assert().Nil(dbEntity.ExpirationTime)

	// When cacheGUID is nil, ID should be kept and no originator info are filled.
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, nil, "client1", 1)
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().Equal(dbEntity.ID, *pbEntity.IdString)
	suite.Assert().Nil(dbEntity.OriginatorCacheGUID)
	suite.Assert().Nil(dbEntity.OriginatorClientItemID)
	suite.Assert().Nil(dbEntity.ExpirationTime)

	// Check that when updating from a previous version with guid, ID will not be
	// replaced.
	pbEntity.Version = aws.Int64(1)
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().Equal(dbEntity.ID, *pbEntity.IdString)
	suite.Assert().Nil(dbEntity.Deleted, "Deleted won't apply its default value for updated entities")
	suite.Assert().Nil(dbEntity.Folder, "Deleted won't apply its default value for updated entities")
	suite.Assert().Nil(dbEntity.ExpirationTime)

	// Empty unique position should be marshalled to nil without error.
	pbEntity.UniquePosition = nil
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Require().NoError(err)
	suite.Assert().Nil(dbEntity.UniquePosition)
	suite.Assert().Nil(dbEntity.ExpirationTime)

	// A history entity should have the client tag hash as the ID,
	// and an expiration time.
	historyEntitySpecific := &sync_pb.EntitySpecifics_History{}
	pbEntity.Specifics = &sync_pb.EntitySpecifics{SpecificsVariant: historyEntitySpecific}
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Require().NoError(err)
	expectedExpirationTime := time.Now().Unix() + datastore.HistoryExpirationIntervalSecs
	suite.Assert().Greater(*dbEntity.ExpirationTime+2, expectedExpirationTime)
	suite.Assert().Less(*dbEntity.ExpirationTime-2, expectedExpirationTime)

	// Empty specifics should report marshal error.
	pbEntity.Specifics = nil
	_, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1", 1)
	suite.Assert().NotNil(err.Error(), "empty specifics should fail")
}

func (suite *SyncEntityTestSuite) TestCreatePBSyncEntity() {
	nigoriSpecific := &sync_pb.NigoriSpecifics{}
	nigoriEntitySpecific := &sync_pb.EntitySpecifics_Nigori{Nigori: nigoriSpecific}
	specifics := &sync_pb.EntitySpecifics{SpecificsVariant: nigoriEntitySpecific}
	specificsBytes, err := proto.Marshal(specifics)
	suite.Require().NoError(err, "Marshal specifics should succeed")

	uniquePosition := &sync_pb.UniquePosition{
		CustomCompressedV1: []byte{1, 2},
	}
	uniquePositionBytes, err := proto.Marshal(uniquePosition)
	suite.Require().NoError(err, "Marshal unique position should succeed")

	id, _ := uuid.NewV7()

	dbEntity := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     id.String(),
		ParentID:               aws.String("parent_id"),
		Version:                aws.Int64(10),
		Mtime:                  aws.Int64(12345678),
		Ctime:                  aws.Int64(12345678),
		Name:                   aws.String("name"),
		NonUniqueName:          aws.String("non_unique_name"),
		ServerDefinedUniqueTag: aws.String("server_tag"),
		ClientDefinedUniqueTag: aws.String("client_tag"),
		Deleted:                aws.Bool(false),
		Folder:                 aws.Bool(false),
		Specifics:              specificsBytes,
		UniquePosition:         uniquePositionBytes,
		DataType:               aws.Int(47745), // nigori type ID
		OriginatorCacheGUID:    aws.String("guid"),
		OriginatorClientItemID: aws.String("client_item_id"),
		DataTypeMtime:          aws.String("47745#12345678"),
	}
	expectedPBEntity := sync_pb.SyncEntity{
		IdString:               &dbEntity.ID,
		ParentIdString:         dbEntity.ParentID,
		Version:                dbEntity.Version,
		Mtime:                  dbEntity.Mtime,
		Ctime:                  dbEntity.Ctime,
		Name:                   dbEntity.Name,
		NonUniqueName:          dbEntity.NonUniqueName,
		ServerDefinedUniqueTag: dbEntity.ServerDefinedUniqueTag,
		ClientTagHash:          dbEntity.ClientDefinedUniqueTag,
		OriginatorCacheGuid:    dbEntity.OriginatorCacheGUID,
		OriginatorClientItemId: dbEntity.OriginatorClientItemID,
		Deleted:                dbEntity.Deleted,
		Folder:                 dbEntity.Folder,
		Specifics:              specifics,
		UniquePosition:         uniquePosition,
	}

	pbEntity, err := datastore.CreatePBSyncEntity(&dbEntity)
	suite.Require().NoError(err, "CreatePBSyncEntity should succeed")

	// Marshal to json to ignore protobuf internal fields when checking equality.
	s1, err := json.Marshal(pbEntity)
	suite.Require().NoError(err, "json.Marshal should succeed")
	s2, err := json.Marshal(&expectedPBEntity)
	suite.Require().NoError(err, "json.Marshal should succeed")
	suite.Assert().Equal(s1, s2)

	// Ensure ID is the client tag for history items
	expectedPBEntity.IdString = expectedPBEntity.ClientTagHash
	*dbEntity.DataType = datastore.HistoryTypeID
	pbEntity, err = datastore.CreatePBSyncEntity(&dbEntity)
	suite.Require().NoError(err, "CreatePBSyncEntity should succeed")

	// Marshal to json to ignore protobuf internal fields when checking equality.
	s1, err = json.Marshal(pbEntity)
	suite.Require().NoError(err, "json.Marshal should succeed")
	s2, err = json.Marshal(&expectedPBEntity)
	suite.Require().NoError(err, "json.Marshal should succeed")
	suite.Assert().Equal(s1, s2)

	// Nil UniquePosition should be unmarshalled as nil without error.
	dbEntity.UniquePosition = nil
	pbEntity, err = datastore.CreatePBSyncEntity(&dbEntity)
	suite.Require().NoError(err, "CreatePBSyncEntity should succeed")
	suite.Assert().Nil(pbEntity.UniquePosition)

	// Nil Specifics should be unmarshalled as nil without error.
	dbEntity.Specifics = nil
	pbEntity, err = datastore.CreatePBSyncEntity(&dbEntity)
	suite.Require().NoError(err, "CreatePBSyncEntity should succeed")
	suite.Assert().Nil(pbEntity.Specifics)
}

func TestSyncEntityTestSuite(t *testing.T) {
	suite.Run(t, new(SyncEntityTestSuite))
}
