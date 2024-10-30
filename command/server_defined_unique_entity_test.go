package command_test

import (
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/stretchr/testify/suite"
)

type ServerDefinedUniqueEntityTestSuite struct {
	suite.Suite
	sqlDB    *datastore.SQLDB
	dynamoDB *datastore.Dynamo
}

type SyncAttrs struct {
	ClientID string
	Name     *string
	DataType *int
	ParentID *string
	Version  *int64
	Deleted  *bool
	Folder   *bool
}

func (suite *ServerDefinedUniqueEntityTestSuite) SetupSuite() {
	datastore.Table = "client-entity-test-command"
	var err error
	suite.dynamoDB, err = datastore.NewDynamo(true)
	suite.Require().NoError(err, "Failed to get dynamoDB session")
	suite.sqlDB, err = datastore.NewSQLDB(true)
	suite.Require().NoError(err, "Failed to get SQL DB session")
}

func (suite *ServerDefinedUniqueEntityTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetDynamoTable(suite.dynamoDB), "Failed to reset Dynamo table")
	suite.Require().NoError(
		datastoretest.ResetSQLTables(suite.sqlDB), "Failed to reset SQL tables")
}

func (suite *ServerDefinedUniqueEntityTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamoDB), "Failed to delete table")
}

func (suite *ServerDefinedUniqueEntityTestSuite) TestInsertServerDefinedUniqueEntities() {
	dbHelpers, err := command.NewDBHelpers(suite.dynamoDB, suite.sqlDB, "client1", nil, false)
	suite.Require().NoError(err, "NewDBHelpers should succeed")
	defer dbHelpers.Trx.Rollback()

	suite.Require().NoError(
		dbHelpers.InsertServerDefinedUniqueEntities(),
		"InsertServerDefinedUniqueEntities should succeed")
	suite.Require().NoError(
		dbHelpers.InsertServerDefinedUniqueEntities(),
		"InsertServerDefinedUniqueEntities again for a same client should succeed")

	expectedSyncAttrsMap := map[string]*SyncAttrs{
		command.NigoriTag: {
			ClientID: "client1",
			Name:     aws.String(command.NigoriName),
			DataType: aws.Int(int(nigoriType)),
			ParentID: aws.String("0"),
			Version:  aws.Int64(1),
			Deleted:  aws.Bool(false),
			Folder:   aws.Bool(true),
		},
		command.BookmarksTag: {
			ClientID: "client1",
			Name:     aws.String(command.BookmarksName),
			DataType: aws.Int(int(bookmarkType)),
			ParentID: aws.String("0"),
			Version:  aws.Int64(1),
			Deleted:  aws.Bool(false),
			Folder:   aws.Bool(true),
		},
		command.OtherBookmarksTag: {
			ClientID: "client1",
			Name:     aws.String(command.OtherBookmarksName),
			DataType: aws.Int(int(bookmarkType)),
			Version:  aws.Int64(1),
			Deleted:  aws.Bool(false),
			Folder:   aws.Bool(true),
		},
		command.SyncedBookmarksTag: {
			ClientID: "client1",
			Name:     aws.String(command.SyncedBookmarksName),
			DataType: aws.Int(int(bookmarkType)),
			Version:  aws.Int64(1),
			Deleted:  aws.Bool(false),
			Folder:   aws.Bool(true),
		},
		command.BookmarkBarTag: {
			ClientID: "client1",
			Name:     aws.String(command.BookmarkBarName),
			DataType: aws.Int(int(bookmarkType)),
			Version:  aws.Int64(1),
			Deleted:  aws.Bool(false),
			Folder:   aws.Bool(true),
		},
	}

	// Tag items should be inserted for each server tags.
	expectedTagItems := []datastore.ServerClientUniqueTagItem{}
	for key := range expectedSyncAttrsMap {
		expectedTagItems = append(expectedTagItems,
			datastore.ServerClientUniqueTagItem{ClientID: "client1", ID: "Server#" + key})
	}
	tagItems, err := datastoretest.ScanTagItems(suite.dynamoDB)
	suite.Require().NoError(err, "ScanTagItems should succeed")

	// Check that Ctime and Mtime have been set, reset to zero value for subsequent
	// tests
	for i := 0; i < len(tagItems); i++ {
		suite.Assert().NotNil(tagItems[i].Ctime)
		suite.Assert().NotNil(tagItems[i].Mtime)

		tagItems[i].Ctime = nil
		tagItems[i].Mtime = nil
	}

	sort.Sort(datastore.TagItemByClientIDID(tagItems))
	sort.Sort(datastore.TagItemByClientIDID(expectedTagItems))
	suite.Assert().Equal(tagItems, expectedTagItems)

	syncItems, err := datastoretest.ScanSyncEntities(suite.dynamoDB)
	suite.Require().NoError(err, "ScanSyncEntities should succeed")

	// Find bookmark root folder to update parentID of its subfolders.
	var bookmarksRootID string
	for _, item := range syncItems {
		if *item.ServerDefinedUniqueTag == command.BookmarksTag {
			bookmarksRootID = item.ID
			expectedSyncAttrsMap[command.OtherBookmarksTag].ParentID = &bookmarksRootID
			expectedSyncAttrsMap[command.SyncedBookmarksTag].ParentID = &bookmarksRootID
			expectedSyncAttrsMap[command.BookmarkBarTag].ParentID = &bookmarksRootID
			break
		}
	}
	suite.Assert().NotEqual(bookmarksRootID, "", "Cannot find ID of bookmarks root folder")

	// For each item returned by ScanSyncEntities, make sure it is in the map and
	// its value is matched, then remove it from the map.
	for _, item := range syncItems {
		syncAttrs := SyncAttrs{
			ClientID: item.ClientID,
			Name:     item.Name,
			DataType: item.DataType,
			ParentID: item.ParentID,
			Version:  item.Version,
			Deleted:  item.Deleted,
			Folder:   item.Folder,
		}

		suite.Assert().NotNil(item.ServerDefinedUniqueTag)
		suite.Assert().Equal(syncAttrs, *expectedSyncAttrsMap[*item.ServerDefinedUniqueTag])
		delete(expectedSyncAttrsMap, *item.ServerDefinedUniqueTag)
	}
	suite.Assert().Equal(0, len(expectedSyncAttrsMap))

	suite.Require().NoError(dbHelpers.Trx.Commit(), "Transaction commit should succeed")

	dbHelpers, err = command.NewDBHelpers(suite.dynamoDB, suite.sqlDB, "client2", nil, false)
	suite.Require().NoError(err, "NewDBHelpers should succeed")
	defer dbHelpers.Trx.Rollback()

	suite.Require().NoError(
		dbHelpers.InsertServerDefinedUniqueEntities(),
		"InsertServerDefinedUniqueEntities should succeed for another client")
}

func TestServerDefinedUniqueEntityTestSuite(t *testing.T) {
	suite.Run(t, new(ServerDefinedUniqueEntityTestSuite))
}
