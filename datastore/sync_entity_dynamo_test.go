package datastore_test

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/utils"
	"github.com/stretchr/testify/suite"
)

type SyncEntityDynamoTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *SyncEntityDynamoTestSuite) SetupSuite() {
	datastore.Table = "client-entity-test-datastore"
	var err error
	suite.dynamo, err = datastore.NewDynamo(true)
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *SyncEntityDynamoTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetDynamoTable(suite.dynamo), "Failed to reset table")
}

func (suite *SyncEntityDynamoTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
}

func (suite *SyncEntityDynamoTestSuite) TestInsertSyncEntity() {
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

func (suite *SyncEntityDynamoTestSuite) TestHasServerDefinedUniqueTag() {
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

func (suite *SyncEntityDynamoTestSuite) TestHasItem() {
	// Insert entity which will be checked later
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
	entity2.ClientID = "client2"
	entity2.ID = "id2"

	_, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity2)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	hasTag, err := suite.dynamo.HasItem("client1", "id1")
	suite.Require().NoError(err, "HasItem should succeed")
	suite.Assert().Equal(hasTag, true)

	hasTag, err = suite.dynamo.HasItem("client2", "id2")
	suite.Require().NoError(err, "HasItem should succeed")
	suite.Assert().Equal(hasTag, true)

	hasTag, err = suite.dynamo.HasItem("client2", "id3")
	suite.Require().NoError(err, "HasItem should succeed")
	suite.Assert().Equal(hasTag, false)

	hasTag, err = suite.dynamo.HasItem("client3", "id2")
	suite.Require().NoError(err, "HasItem should succeed")
	suite.Assert().Equal(hasTag, false)
}

func (suite *SyncEntityDynamoTestSuite) TestInsertSyncEntitiesWithServerTags() {
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

func (suite *SyncEntityDynamoTestSuite) TestUpdateSyncEntity_Basic() {
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
	suite.Assert().True(deleted, "Successful update should result in delete")

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
	suite.Assert().False(deleted, "Successful update should not result in delete")

	// Update with nil Folder and Deleted
	updateEntity3 := updateEntity1
	updateEntity3.ID = "id3"
	updateEntity3.Folder = nil
	updateEntity3.Deleted = nil
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity3, *entity3.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().False(deleted, "Successful update should not result in delete")
	// Reset these back to false because they will be the expected value in DB.
	updateEntity3.Folder = aws.Bool(false)
	updateEntity3.Deleted = aws.Bool(false)

	// Update entity again with the wrong old version as (version mismatch)
	// should return false.
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity2, 12345678)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Update with the same version should return conflict")
	suite.Assert().False(deleted, "Successful update should not result in delete")
	// suite.Assert().False(deleted, "Successful update should not result in delete")

	// Check sync entities are updated correctly in DB.
	syncItems, err = datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{updateEntity1, updateEntity2, updateEntity3})
}

func (suite *SyncEntityDynamoTestSuite) TestUpdateSyncEntity_HistoryType() {
	// Insert a history item
	entity1 := datastore.SyncEntity{
		ClientID:               "client1",
		ID:                     "id1",
		Version:                aws.Int64(1),
		ClientDefinedUniqueTag: aws.String("client_tag1"),
		Ctime:                  aws.Int64(12345678),
		Mtime:                  aws.Int64(12345678),
		DataType:               aws.Int(963985),
		Folder:                 aws.Bool(false),
		Deleted:                aws.Bool(false),
		DataTypeMtime:          aws.String("123#12345678"),
		Specifics:              []byte{1, 2},
	}
	conflict, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful insert should not have conflict")

	updateEntity1 := entity1
	updateEntity1.Version = aws.Int64(2)
	updateEntity1.Folder = aws.Bool(true)
	updateEntity1.Mtime = aws.Int64(24242424)
	conflict, deleted, err := suite.dynamo.UpdateSyncEntity(&updateEntity1, 1)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().False(deleted, "Successful update should not result in delete")

	// should still succeed with the same version number,
	// since the version number should be ignored
	updateEntity2 := updateEntity1
	updateEntity2.Mtime = aws.Int64(42424242)
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity2, 1)
	suite.Require().NoError(err, "UpdateSyncEntity should not return an error")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().False(deleted, "Successful update should not result in delete")

	updateEntity3 := entity1
	updateEntity3.Deleted = aws.Bool(true)

	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity3, 1)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().True(deleted, "Successful update should result in delete")

	syncItems, err := datastoretest.ScanSyncEntities(suite.dynamo)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")
	updateEntity3.ID = *updateEntity3.ClientDefinedUniqueTag
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{updateEntity3})
}

func (suite *SyncEntityDynamoTestSuite) TestUpdateSyncEntity_ReuseClientTag() {
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
	suite.Assert().False(deleted, "Successful update should not result in delete")

	// Soft-delete the item with wrong version should get conflict.
	updateEntity1.Deleted = aws.Bool(true)
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity1, *entity1.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Version mismatched update should have conflict")
	suite.Assert().False(deleted, "Successful update should not result in delete")

	// Soft-delete the item with matched version.
	updateEntity1.Version = aws.Int64(34567890)
	conflict, deleted, err = suite.dynamo.UpdateSyncEntity(&updateEntity1, 23456789)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Successful update should not have conflict")
	suite.Assert().True(deleted, "Successful update should result in delete")

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

func (suite *SyncEntityDynamoTestSuite) TestGetUpdatesForType() {
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

	// non-expired item
	entity4 := entity2
	entity4.ClientID = "client2"
	entity4.ID = "id4"
	entity4.ExpirationTime = aws.Int64(time.Now().Unix() + 300)

	// expired item
	entity5 := entity2
	entity5.ClientID = "client2"
	entity5.ID = "id5"
	entity5.ExpirationTime = aws.Int64(time.Now().Unix() - 300)

	_, err := suite.dynamo.InsertSyncEntity(&entity1)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity2)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity3)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity4)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	_, err = suite.dynamo.InsertSyncEntity(&entity5)
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	// Get all updates for type 123 and client1 using token = 0.
	var token int64
	hasChangesRemaining, syncItems, err := suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client1", 100, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1, entity2})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for type 124 and client1 using token = 0.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(124, &token, nil, true, "client1", 100, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity3})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for type 123 and client2 using token = 0.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client2", 100, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity4})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for type 124 and client2 using token = 0.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(124, &token, nil, true, "client2", 100, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(len(syncItems), 0)
	suite.Assert().False(hasChangesRemaining)

	// Test maxSize will limit the return entries size, and hasChangesRemaining
	// should be true when there are more updates available in the DB.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client1", 1, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1})
	suite.Assert().True(hasChangesRemaining)

	// Test when num of query items equal to the limit, hasChangesRemaining should
	// be true.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client1", 2, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity1, entity2})
	suite.Assert().True(hasChangesRemaining)

	// Test fetchFolders will remove folder items if false
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, false, "client1", 100, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity2})
	suite.Assert().False(hasChangesRemaining)

	// Get all updates for a type for a client using mtime of one item as token.
	token = 12345678
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client1", 100, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, []datastore.SyncEntity{entity2})
	suite.Assert().False(hasChangesRemaining)

	// Test batch is working correctly for over 100 items
	err = datastoretest.ResetDynamoTable(suite.dynamo)
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
	token = 0
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client1", 300, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	sort.Sort(datastore.SyncEntityByMtime(expectedSyncItems))
	suite.Assert().Equal(syncItems, expectedSyncItems)
	suite.Assert().False(hasChangesRemaining)

	// Test that when maxGUBatchSize is smaller than total updates, the first n
	// items ordered by Mtime should be returned.
	hasChangesRemaining, syncItems, err = suite.dynamo.GetUpdatesForType(123, &token, nil, true, "client1", 200, true)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().Equal(syncItems, expectedSyncItems[0:200])
	suite.Assert().True(hasChangesRemaining)
}

func (suite *SyncEntityDynamoTestSuite) TestDisableSyncChain() {
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

func (suite *SyncEntityDynamoTestSuite) TestIsSyncChainDisabled() {
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

func (suite *SyncEntityDynamoTestSuite) TestClearServerData() {
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

func TestSyncEntityDynamoTestSuite(t *testing.T) {
	suite.Run(t, new(SyncEntityDynamoTestSuite))
}
