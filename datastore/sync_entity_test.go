package datastore_test

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/brave/go-sync/utils"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
)

type SyncEntityTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *SyncEntityTestSuite) SetupSuite() {
	datastore.Table = "client-entity-test-datastore"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *SyncEntityTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *SyncEntityTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
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

func (suite *SyncEntityTestSuite) TestInsertSyncEntity() {
	entity1 := datastore.SyncEntity{
		ClientID:      "client1",
		ID:            "id1",
		Version:       aws.Int64(1),
		Ctime:         aws.Int64(12345678),
		Mtime:         aws.Int64(12345678),
		DataType:      aws.Int(123),
		Folder:        aws.Bool(false),
		Deleted:       aws.Bool(false),
		DataTypeMtime: aws.String("123#12345678"),
	}
	entity2 := entity1
	entity2.ID = "id2"
	_, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity2)
	suite.Require().NoError(err, "InsertSyncEntity with other ID should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().Error(err, "InsertSyncEntity with the same ClientID and ID should fail")

	// Each InsertSyncEntity without client tag should result in one sync item saved.
	tagItems, err := datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(
		0, len(tagItems), "Insert without client tag should not insert tag items")
	syncItems, err := datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1, entity2})

	// Insert entity with client tag should result in one sync item and one tag
	// item saved.
	entity3 := entity1
	entity3.ID = "id3"
	entity3.ClientDefinedUniqueTag = aws.String("tag1")
	_, err = suite.dynamo.InsertSyncEntity(&entity3)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	// Insert entity with different tag for same ClientID should succeed.
	entity4 := entity3
	entity4.ID = "id4"
	entity4.ClientDefinedUniqueTag = aws.String("tag2")
	_, err = suite.dynamo.InsertSyncEntity(&entity4)
	suite.Require().NoError(err, "InsertSyncEntity with different server tag should succeed")

	// Insert entity with the same client tag and ClientID should fail with conflict.
	entity4Copy := entity4
	entity4Copy.ID = "id4_copy"
	conflict, err := suite.dynamo.InsertSyncEntity(&entity4Copy)
	suite.Require().Error(err, "InsertSyncEntity with the same client tag and ClientID should fail")
	suite.Assert().True(conflict, "Return conflict for duplicate client tag")

	// Insert entity with the same client tag for other client should not fail.
	entity5 := entity3
	entity5.ClientID = "client2"
	entity5.ID = "id5"
	_, err = suite.dynamo.InsertSyncEntity(&entity5)
	suite.Require().NoError(err,
		"InsertSyncEntity with the same client tag for another client should succeed")

	// Check sync items are saved for entity1, entity2, entity3, entity4, entity5.
	syncItems, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	expectedSyncItems := []datastore.SyncEntity{entity1, entity2, entity3, entity4, entity5}
	sort.Sort(datastore.SyncEntityByClientIDID(syncItems))
	suite.Assert().Equal(syncItems, expectedSyncItems)

	// Check tag items should be saved for entity3, entity4, entity5.
	tagItems, err = datastoretest.ScanTagItems(suite.dynamo)

	// Check that Ctime and Mtime have been set, reset to zero value for subsequent
	// tests
	for i := 0; i < len(tagItems); i++ {
		suite.Assert().NotNil(tagItems[i].Ctime)
		suite.Assert().NotNil(tagItems[i].Mtime)

		tagItems[i].Ctime = nil
		tagItems[i].Mtime = nil
	}

	suite.Require().NoError(err, "ScanTagItems should succeed")
	expectedTagItems := []datastore.ServerClientUniqueTagItem{
		{ClientID: "client1", ID: "Client#tag1"},
		{ClientID: "client1", ID: "Client#tag2"},
		{ClientID: "client2", ID: "Client#tag1"},
	}
	sort.Sort(datastore.TagItemByClientIDID(tagItems))
	suite.Assert().Equal(expectedTagItems, tagItems)
}

func (suite *SyncEntityTestSuite) TestHasServerDefinedUniqueTag() {
	// Insert entities with server tags using InsertSyncEntitiesWithServerTags.
	tag1 := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     "id1",
		Version:                aws.Int64(1),
		Ctime:                  aws.Int64(12345678),
		Mtime:                  aws.Int64(12345678),
		DataType:               aws.Int(123),
		Folder:                 aws.Bool(true),
		Deleted:                aws.Bool(false),
		DataTypeMtime:          aws.String("123#12345678"),
		ServerDefinedUniqueTag: aws.String("tag1"),
	}
	tag2 := tag1
	tag2.ClientID = "client2"
	tag2.ID = "id2"
	tag2.ServerDefinedUniqueTag = aws.String("tag2")
	entities := []*datastore.SyncEntity{&tag1, &tag2}

	err := suite.dynamo.InsertSyncEntitiesWithServerTags(entities)
	suite.Require().NoError(err, "Insert sync entities should succeed")

	hasTag, err := suite.dynamo.HasServerDefinedUniqueTag("client1", "tag1")
	suite.Require().NoError(err, "HasServerDefinedUniqueTag should succeed")
	suite.Assert().Equal(hasTag, true)

	hasTag, err = suite.dynamo.HasServerDefinedUniqueTag("client1", "tag2")
	suite.Require().NoError(err, "HasServerDefinedUniqueTag should succeed")
	suite.Assert().Equal(hasTag, false)

	hasTag, err = suite.dynamo.HasServerDefinedUniqueTag("client2", "tag1")
	suite.Require().NoError(err, "HasServerDefinedUniqueTag should succeed")
	suite.Assert().Equal(hasTag, false)

	hasTag, err = suite.dynamo.HasServerDefinedUniqueTag("client2", "tag2")
	suite.Require().NoError(err, "HasServerDefinedUniqueTag should succeed")
	suite.Assert().Equal(hasTag, true)
}

func (suite *SyncEntityTestSuite) TestInsertSyncEntitiesWithServerTags() {
	// Insert with same ClientID and server tag would fail.
	entity1 := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     "id1",
		Version:                aws.Int64(1),
		Ctime:                  aws.Int64(12345678),
		Mtime:                  aws.Int64(12345678),
		DataType:               aws.Int(123),
		Folder:                 aws.Bool(false),
		Deleted:                aws.Bool(false),
		DataTypeMtime:          aws.String("123#12345678"),
		ServerDefinedUniqueTag: aws.String("tag1"),
	}
	entity2 := entity1
	entity2.ID = "id2"
	entities := []*datastore.SyncEntity{&entity1, &entity2}
	suite.Require().Error(
		suite.dynamo.InsertSyncEntitiesWithServerTags(entities),
		"Insert with same ClientID and server tag would fail")

	// Check nothing is written to DB when it fails.
	syncItems, err := datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(0, len(syncItems), "No items should be written if fail")
	tagItems, err := datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(0, len(tagItems), "No items should be written if fail")

	entity2.ServerDefinedUniqueTag = aws.String("tag2")
	entity3 := entity1
	entity3.ClientID = "client2"
	entity3.ID = "id3"
	entities = []*datastore.SyncEntity{&entity1, &entity2, &entity3}
	suite.Require().NoError(
		suite.dynamo.InsertSyncEntitiesWithServerTags(entities),
		"InsertSyncEntitiesWithServerTags should succeed")

	// Scan DB and check all items are saved
	syncItems, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	expectedSyncItems := []datastore.SyncEntity{entity1, entity2, entity3}
	sort.Sort(datastore.SyncEntityByClientIDID(syncItems))
	suite.Assert().Equal(syncItems, expectedSyncItems)
	tagItems, err = datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")

	// Check that Ctime and Mtime have been set, reset to zero value for subsequent
	// tests
	for i := 0; i < len(tagItems); i++ {
		suite.Assert().NotNil(tagItems[i].Ctime)
		suite.Assert().NotNil(tagItems[i].Mtime)

		tagItems[i].Ctime = nil
		tagItems[i].Mtime = nil
	}

	expectedTagItems := []datastore.ServerClientUniqueTagItem{
		{ClientID: "client1", ID: "Server#tag1"},
		{ClientID: "client1", ID: "Server#tag2"},
		{ClientID: "client2", ID: "Server#tag1"},
	}
	sort.Sort(datastore.TagItemByClientIDID(tagItems))
	suite.Assert().Equal(expectedTagItems, tagItems)
}

func (suite *SyncEntityTestSuite) TestUpdateSyncEntity_Basic() {
	// Insert three new items.
	entity1 := datastore.SyncEntity{
		ClientID:      "client1",
		ID:            "id1",
		Version:       aws.Int64(1),
		Ctime:         aws.Int64(12345678),
		Mtime:         aws.Int64(12345678),
		DataType:      aws.Int(123),
		Folder:        aws.Bool(false),
		Deleted:       aws.Bool(false),
		DataTypeMtime: aws.String("123#12345678"),
		Specifics:     []byte{1, 2},
	}
	entity2 := entity1
	entity2.ID = "id2"
	entity3 := entity1
	entity3.ID = "id3"
	_, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity2)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity3)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	// Check sync entities are inserted correctly in DB.
	syncItems, err := datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1, entity2, entity3})

	// Update without optional fields.
	updateEntity1 := entity1
	updateEntity1.Version = aws.Int64(23456789)
	updateEntity1.Mtime = aws.Int64(23456789)
	updateEntity1.Folder = aws.Bool(true)
	updateEntity1.Deleted = aws.Bool(true)
	updateEntity1.DataTypeMtime = aws.String("123#23456789")
	updateEntity1.Specifics = []byte{3, 4}
	conflict, deleted, err := suite.dynamo.UpdateSyncEntity(&updateEntity1, *entity1.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().True(deleted, "Delete operation should return true")

	// Update with optional fields.
	updateEntity2 := updateEntity1
	updateEntity2.ID = "id2"
	updateEntity2.Deleted = aws.Bool(false)
	updateEntity2.Folder = aws.Bool(false)
	updateEntity2.UniquePosition = []byte{5, 6}
	updateEntity2.ParentID = aws.String("parentID")
	updateEntity2.Name = aws.String("name")
	updateEntity2.NonUniqueName = aws.String("non_unique_name")
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity2, *entity2.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().False(deleted, "Non-delete operation should return false")

	// Update with nil Folder and Deleted
	updateEntity3 := updateEntity1
	updateEntity3.ID = "id3"
	updateEntity3.Folder = nil
	updateEntity3.Deleted = nil
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity3, *entity3.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().False(deleted, "Non-delete operation should return false")
	// Reset these back to false because they will be the expected value in DB.
	updateEntity3.Folder = aws.Bool(false)
	updateEntity3.Deleted = aws.Bool(false)

	// Update entity again with the wrong old version as (version mismatch)
	// should return false.
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity2, 12345678)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Update with the same version should return conflict")
	suite.Assert().False(deleted, "Conflict operation should return false for delete")

	// Check sync entities are updated correctly in DB.
	syncItems, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{updateEntity1, updateEntity2, updateEntity3})
}

func (suite *SyncEntityTestSuite) TestUpdateSyncEntity_ReuseClientTag() {
	// Insert an item with client tag.
	entity1 := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     "id1",
		Version:                aws.Int64(1),
		ClientDefinedUniqueTag: aws.String("client_tag"),
		Ctime:                  aws.Int64(12345678),
		Mtime:                  aws.Int64(12345678),
		DataType:               aws.Int(123),
		Folder:                 aws.Bool(false),
		Deleted:                aws.Bool(false),
		DataTypeMtime:          aws.String("123#12345678"),
		Specifics:              []byte{1, 2},
	}
	conflict, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful insert should not have conflict")

	// Check a tag item is inserted.
	tagItems, err := datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(1, len(tagItems), "Tag item should be inserted")

	// Update it to version 23456789.
	updateEntity1 := entity1
	updateEntity1.Version = aws.Int64(23456789)
	updateEntity1.Mtime = aws.Int64(23456789)
	updateEntity1.Folder = aws.Bool(true)
	updateEntity1.DataTypeMtime = aws.String("123#23456789")
	updateEntity1.Specifics = []byte{3, 4}
	conflict, deleted, err := suite.dynamo.UpdateSyncEntity(&updateEntity1, *entity1.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().False(deleted, "Non-delete operation should return false")

	// Soft-delete the item with wrong version should get conflict.
	updateEntity1.Deleted = aws.Bool(true)
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity1, *entity1.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Version mismatched update should have conflict")
	suite.Assert().False(deleted, "Failed delete operation should return false")

	// Soft-delete the item with matched version.
	updateEntity1.Version = aws.Int64(34567890)
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity1, 23456789)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().True(deleted, "Delete operation should return true")

	// Check tag item is deleted.
	tagItems, err = datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(0, len(tagItems), "Tag item should be deleted")

	// Insert another item with the same client tag again.
	entity2 := entity1
	entity2.ID = "id2"
	conflict, err = suite.dynamo.InsertSyncEntity(&entity2)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful insert should not have conflict")

	// Check a tag item is inserted.
	tagItems, err = datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(1, len(tagItems), "Tag item should be inserted")
}

func (suite *SyncEntityTestSuite) TestGetUpdatesForType() {
	// Insert items for testing.
	entity1 := datastore.SyncEntity{
		ClientID:      "client1",
		ID:            "id1",
		Version:       aws.Int64(1),
		Ctime:         aws.Int64(12345678),
		Mtime:         aws.Int64(12345678),
		DataType:      aws.Int(123),
		Folder:        aws.Bool(true),
		Deleted:       aws.Bool(false),
		DataTypeMtime: aws.String("123#12345678"),
		Specifics:     []byte{1, 2},
	}

	entity2 := entity1
	entity2.ID = "id2"
	entity2.Folder = aws.Bool(false)
	entity2.Mtime = aws.Int64(12345679)
	entity2.DataTypeMtime = aws.String("123#12345679")

	entity3 := entity2
	entity3.ID = "id3"
	entity3.DataType = aws.Int(124)
	entity3.DataTypeMtime = aws.String("124#12345679")

	entity4 := entity2
	entity4.ClientID = "client2"
	entity4.ID = "id4"

	_, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity2)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity3)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity4)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	// Get all updates for type 123 and client1 using token = 0.
	hasChangesRemaining, syncItems, err := suite.dynamo.GetUpdatesForType(123, 0, true, "client1", 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1, entity2})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for type 124 and client1 using token = 0.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(124, 0, true, "client1", 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity3})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for type 123 and client2 using token = 0.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 0, true, "client2", 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity4})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for type 124 and client2 using token = 0.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(124, 0, true, "client2", 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(len(syncItems), 0)
	suite.Assert().False(hasChangesRemaining)

	// Test maxSize will limit the return entries size, and hasChangesRemaining
	// should be true when there are more updates available in the DB.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 0, true, "client1", 1)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1})
	suite.Assert().True(hasChangesRemaining)

	// Test when num of query items equal to the limit, hasChangesRemaining should
	// be true.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 0, true, "client1", 2)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1, entity2})
	suite.Assert().True(hasChangesRemaining)

	// Test fetchFolders will remove folder items if false
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 0, false, "client1", 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity2})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for a type for a client using mtime of one item as token.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 12345678, true, "client1", 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity2})
	suite.Assert().False(hasChangesRemaining)

	// Test batch is working correctly for over 100 items
	err = datastoretest.ResetTable(suite.dynamo)
	suite.Require().NoError(err, "Failed to reset table")

	expectedSyncItems := []datastore.SyncEntity{}
	entity1 = datastore.SyncEntity{
		ClientID:  "client1",
		Version:   aws.Int64(1),
		Ctime:     aws.Int64(12345678),
		DataType:  aws.Int(123),
		Folder:    aws.Bool(false),
		Deleted:   aws.Bool(false),
		Specifics: []byte{1, 2},
	}

	mtime := utils.UnixMilli(time.Now())
	for i := 1; i <= 250; i++ {
		mtime = mtime + 1
		entity := entity1
		entity.ID = "id" + strconv.Itoa(i)
		entity.Mtime = aws.Int64(mtime)
		entity.DataTypeMtime = aws.String("123#" + strconv.FormatInt(*entity.Mtime, 10))
		_, err := suite.dynamo.InsertSyncEntity(&entity)
		suite.Require().NoError(err, "InsertSyncEntity should succeed")
		expectedSyncItems = append(expectedSyncItems, entity)
	}

	// All items should be returned and sorted by Mtime.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 0, true, "client1", 300)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	sort.Sort(datastore.SyncEntityByMtime(expectedSyncItems))
	suite.Assert().Equal(syncItems, expectedSyncItems)
	suite.Assert().False(hasChangesRemaining)

	// Test that when maxGUBatchSize is smaller than total updates, the first n
	// items ordered by Mtime should be returned.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, 0, true, "client1", 200)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, expectedSyncItems[0:200])
	suite.Assert().True(hasChangesRemaining)
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
	expectedDBEntity := datastore.SyncEntity{
		ClientID:               "client1",
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
	}

	dbEntity, err := datastore.CreateDBSyncEntity(&pbEntity, guid, "client1")
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")

	// Check ID is replaced with a server-generated ID.
	suite.Assert().NotEqual(
		dbEntity.ID, *pbEntity.IdString,
		"ID should be a server-generated ID and not equal to the passed IdString")
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

	pbEntity.Deleted = nil
	pbEntity.Folder = nil
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1")
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().False(*dbEntity.Deleted, "Default value should be set for Deleted for new entities")
	suite.Assert().False(*dbEntity.Folder, "Default value should be set for Deleted for new entities")

	// Check the case when Ctime and Mtime are provided by the client.
	pbEntity.Ctime = aws.Int64(12345678)
	pbEntity.Mtime = aws.Int64(12345678)
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1")
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().Equal(*dbEntity.Ctime, *pbEntity.Ctime, "Client's Ctime should be respected")
	suite.Assert().NotEqual(*dbEntity.Mtime, *pbEntity.Mtime, "Client's Mtime should be replaced")

	// When cacheGUID is nil, ID should be kept and no originator info are filled.
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, nil, "client1")
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().Equal(dbEntity.ID, *pbEntity.IdString)
	suite.Assert().Nil(dbEntity.OriginatorCacheGUID)
	suite.Assert().Nil(dbEntity.OriginatorClientItemID)

	// Check that when updating from a previous version with guid, ID will not be
	// replaced.
	pbEntity.Version = aws.Int64(1)
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1")
	suite.Require().NoError(err, "CreateDBSyncEntity should succeed")
	suite.Assert().Equal(dbEntity.ID, *pbEntity.IdString)
	suite.Assert().Nil(dbEntity.Deleted, "Deleted won't apply its default value for updated entities")
	suite.Assert().Nil(dbEntity.Folder, "Deleted won't apply its default value for updated entities")

	// Empty unique position should be marshalled to nil without error.
	pbEntity.UniquePosition = nil
	dbEntity, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1")
	suite.Require().NoError(err)
	suite.Assert().Nil(dbEntity.UniquePosition)

	// Empty specifics should report marshal error.
	pbEntity.Specifics = nil
	_, err = datastore.CreateDBSyncEntity(&pbEntity, guid, "client1")
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

	dbEntity := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     "id1",
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

func (suite *SyncEntityTestSuite) TestDisableSyncChain() {
	clientID := "client1"
	id := "disabled_chain"
	err := suite.dynamo.DisableSyncChain(clientID)
	suite.Require().NoError(err, "DisableSyncChain should succeed")
	e, err := datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(1, len(e))
	suite.Assert().Equal(clientID, e[0].ClientID)
	suite.Assert().Equal(id, e[0].ID)
}

func (suite *SyncEntityTestSuite) TestIsSyncChainDisabled() {
	clientID := "client1"

	disabled, err := suite.dynamo.IsSyncChainDisabled(clientID)
	suite.Require().NoError(err, "IsSyncChainDisabled should succeed")
	suite.Assert().Equal(false, disabled)

	err = suite.dynamo.DisableSyncChain(clientID)
	suite.Require().NoError(err, "DisableSyncChain should succeed")
	disabled, err = suite.dynamo.IsSyncChainDisabled(clientID)
	suite.Require().NoError(err, "IsSyncChainDisabled should succeed")
	suite.Assert().Equal(true, disabled)
}

func (suite *SyncEntityTestSuite) TestClearServerData() {
	// Test clear sync entities
	entity := datastore.SyncEntity{
		ClientID:      "client1",
		ID:            "id1",
		Version:       aws.Int64(1),
		Ctime:         aws.Int64(12345678),
		Mtime:         aws.Int64(12345678),
		DataType:      aws.Int(123),
		Folder:        aws.Bool(false),
		Deleted:       aws.Bool(false),
		DataTypeMtime: aws.String("123#12345678"),
	}
	_, err := suite.dynamo.InsertSyncEntity(&entity)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	e, err := datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(1, len(e))

	e, err = suite.dynamo.ClearServerData(entity.ClientID)
	suite.Require().NoError(err, "ClearServerData should succeed")
	suite.Assert().Equal(1, len(e))

	e, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(0, len(e))

	// Test clear tagged items
	entity1 := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     "id1",
		Version:                aws.Int64(1),
		Ctime:                  aws.Int64(12345678),
		Mtime:                  aws.Int64(12345678),
		DataType:               aws.Int(123),
		Folder:                 aws.Bool(false),
		Deleted:                aws.Bool(false),
		DataTypeMtime:          aws.String("123#12345678"),
		ServerDefinedUniqueTag: aws.String("tag1"),
	}
	entity2 := entity1
	entity2.ID = "id2"
	entity2.ServerDefinedUniqueTag = aws.String("tag2")
	entities := []*datastore.SyncEntity{&entity1, &entity2}
	suite.Require().NoError(
		suite.dynamo.InsertSyncEntitiesWithServerTags(entities),
		"InsertSyncEntitiesWithServerTags should succeed")

	e, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(2, len(e), "No items should be written if fail")

	t, err := datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(2, len(t), "No items should be written if fail")

	e, err = suite.dynamo.ClearServerData(entity.ClientID)
	suite.Require().NoError(err, "ClearServerData should succeed")
	suite.Assert().Equal(4, len(e))

	e, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(0, len(e), "No items should be written if fail")

	t, err = datastoretest.ScanTagItems(suite.dynamo)
	suite.Require().NoError(err, "ScanTagItems should succeed")
	suite.Assert().Equal(0, len(t), "No items should be written if fail")
}

func TestSyncEntityTestSuite(t *testing.T) {
	suite.Run(t, new(SyncEntityTestSuite))
}
