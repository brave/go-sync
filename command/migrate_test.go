package command_test

import (
	"context"
	"math"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/stretchr/testify/suite"
)

type CommandMigrateTestSuite struct {
	suite.Suite
	dynamoDB *datastore.Dynamo
	cache    *cache.Cache
	sqlDB    *datastore.SQLDB
}

func (suite *CommandMigrateTestSuite) SetupSuite() {
	datastore.Table = testDynamoTable
	var err error
	suite.dynamoDB, err = datastore.NewDynamo(true)
	suite.Require().NoError(err, "Failed to get dynamoDB session")

	suite.cache = cache.NewCache(cache.NewRedisClient())
}

type ExpectedCounts struct {
	SQLNigori      int64
	SQLBookmark    int64
	DynamoNigori   int64
	DynamoBookmark int64
}

func (suite *CommandMigrateTestSuite) assertDatastoreCounts(expected ExpectedCounts) {
	sqlNigoriCount, err := getDatastoreCount(true, suite.dynamoDB, suite.sqlDB, []int32{nigoriType})
	suite.Require().NoError(err, "Failed to get SQL nigori count")

	sqlBookmarkCount, err := getDatastoreCount(true, suite.dynamoDB, suite.sqlDB, []int32{bookmarkType})
	suite.Require().NoError(err, "Failed to get SQL bookmark count")

	dynamoNigoriCount, err := getDatastoreCount(false, suite.dynamoDB, suite.sqlDB, []int32{nigoriType})
	suite.Require().NoError(err, "Failed to get DynamoDB nigori count")

	dynamoBookmarkCount, err := getDatastoreCount(false, suite.dynamoDB, suite.sqlDB, []int32{bookmarkType})
	suite.Require().NoError(err, "Failed to get DynamoDB bookmark count")

	suite.Assert().Equal(expected.SQLNigori, sqlNigoriCount, "SQL nigori count mismatch")
	suite.Assert().Equal(expected.SQLBookmark, sqlBookmarkCount, "SQL bookmark count mismatch")
	suite.Assert().Equal(expected.DynamoNigori, dynamoNigoriCount, "DynamoDB nigori count mismatch")
	suite.Assert().Equal(expected.DynamoBookmark, dynamoBookmarkCount, "DynamoDB bookmark count mismatch")
}

func (suite *CommandMigrateTestSuite) assertSQLMigrationStatus(dataType int32, checkForFullMigration bool, shouldExist bool) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM dynamo_migration_statuses
		WHERE data_type = $1`

	if checkForFullMigration {
		query += ` AND earliest_mtime IS NULL`
	}

	err := suite.sqlDB.QueryRow(query, dataType).Scan(&count)

	var expectedCount int
	if shouldExist {
		expectedCount = 1
	}

	suite.Require().NoError(err, "Failed to query dynamo_migration_statuses")
	suite.Assert().Equal(expectedCount, count, "Migration status row count should match")
}

func (suite *CommandMigrateTestSuite) createSQLDB(migrateDataTypes []int32) {
	rollouts := buildRolloutConfigString(migrateDataTypes)
	suite.T().Setenv(datastore.SQLSaveRolloutsEnvKey, rollouts)
	suite.T().Setenv(datastore.SQLMigrateRolloutsEnvKey, rollouts)
	suite.T().Setenv(datastore.SQLMigrateChunkSizeEnvKey, "2")
	suite.T().Setenv(datastore.SQLMigrateUpdateIntervalEnvKey, "1")

	isFirstRun := suite.sqlDB == nil

	var err error
	suite.sqlDB, err = datastore.NewSQLDB(true)
	suite.Require().NoError(err, "Failed to get SQL DB session")

	if isFirstRun {
		suite.Require().NoError(
			datastoretest.ResetSQLTables(suite.sqlDB), "Failed to reset SQL tables")
	}
}

func (suite *CommandMigrateTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetDynamoTable(suite.dynamoDB), "Failed to reset Dynamo table")
}

func (suite *CommandMigrateTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamoDB), "Failed to delete table")
	suite.Require().NoError(
		suite.cache.FlushAll(context.Background()), "Failed to clear cache")
	suite.sqlDB = nil
}

func (suite *CommandMigrateTestSuite) sendMessageAndAssertEmptyResponse(msg *sync_pb.ClientToServerMessage) {
	rsp := &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamoDB, suite.sqlDB, testClientID),
		"HandleClientToServerMessage should succeed")

	suite.Assert().Equal(sync_pb.SyncEnums_SUCCESS, *rsp.ErrorCode, "errorCode should match")
	suite.Assert().NotNil(rsp.GetUpdates)
	suite.Assert().Empty(rsp.GetUpdates.Entries)
}

func (suite *CommandMigrateTestSuite) TestBasicMigrate() {
	suite.createSQLDB([]int32{})
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id3_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id4_nigori", 0, false, getNigoriSpecifics()),
		getCommitEntity("id5_nigori", 0, false, getNigoriSpecifics()),
		getCommitEntity("id6_nigori", 0, false, getNigoriSpecifics()),
		getCommitEntity("id7_nigori", 0, false, getNigoriSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	// Commit and check response.
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamoDB, suite.sqlDB, testClientID),
		"HandleClientToServerMessage should succeed")

	// GetUpdates should return nothing.
	marker := getMarker(MarkerTokens{
		Nigori:   aws.Int64(math.MaxInt64 - 1000),
		Bookmark: aws.Int64(math.MaxInt64 - 1000),
	})

	msg = getClientToServerGUMsg(
		marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)
	suite.sendMessageAndAssertEmptyResponse(msg)

	isSQLEmpty, err := verifyNoDataInOtherDB(false, suite.dynamoDB, suite.sqlDB)
	suite.Require().NoError(err, "Empty database verification should succeed")
	suite.Assert().True(isSQLEmpty, "SQL database should be empty")

	suite.createSQLDB([]int32{nigoriType, bookmarkType})

	suite.sendMessageAndAssertEmptyResponse(msg)
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      2,
		SQLBookmark:    0,
		DynamoNigori:   2,
		DynamoBookmark: 3,
	})
	suite.assertSQLMigrationStatus(bookmarkType, false, false)
	suite.assertSQLMigrationStatus(nigoriType, false, true)

	suite.sendMessageAndAssertEmptyResponse(msg)
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      4,
		SQLBookmark:    0,
		DynamoNigori:   0,
		DynamoBookmark: 3,
	})
	suite.assertSQLMigrationStatus(bookmarkType, false, false)
	suite.assertSQLMigrationStatus(nigoriType, false, true)

	suite.sendMessageAndAssertEmptyResponse(msg)
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      4,
		SQLBookmark:    2,
		DynamoNigori:   0,
		DynamoBookmark: 1,
	})
	suite.assertSQLMigrationStatus(bookmarkType, false, true)
	suite.assertSQLMigrationStatus(nigoriType, true, true)

	suite.sendMessageAndAssertEmptyResponse(msg)
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      4,
		SQLBookmark:    3,
		DynamoNigori:   0,
		DynamoBookmark: 0,
	})
	suite.assertSQLMigrationStatus(bookmarkType, true, true)
	suite.assertSQLMigrationStatus(nigoriType, true, true)

	suite.sendMessageAndAssertEmptyResponse(msg)
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      4,
		SQLBookmark:    3,
		DynamoNigori:   0,
		DynamoBookmark: 0,
	})
	// fully migrated
	suite.assertSQLMigrationStatus(bookmarkType, true, true)
	suite.assertSQLMigrationStatus(nigoriType, true, true)
}

func (suite *CommandMigrateTestSuite) TestBookmarkOnlyMigration() {
	suite.createSQLDB([]int32{})
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id3_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id4_nigori", 0, false, getNigoriSpecifics()),
		getCommitEntity("id5_nigori", 0, false, getNigoriSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	// Commit initial entities
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamoDB, suite.sqlDB, testClientID),
		"HandleClientToServerMessage should succeed")

	// Enable migration for bookmarks only
	suite.createSQLDB([]int32{bookmarkType})

	// GetUpdates message
	marker := getMarker(MarkerTokens{
		Nigori:   aws.Int64(math.MaxInt64 - 1000),
		Bookmark: aws.Int64(math.MaxInt64 - 1000),
	})
	msg = getClientToServerGUMsg(marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)

	// Initial counts
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      0,
		SQLBookmark:    0,
		DynamoNigori:   2,
		DynamoBookmark: 3,
	})

	// Migrate bookmarks
	for i := 0; i < 4; i++ {
		suite.sendMessageAndAssertEmptyResponse(msg)
		if i == 0 {
			suite.assertSQLMigrationStatus(bookmarkType, false, true)
		}
	}

	// Final counts
	suite.assertDatastoreCounts(ExpectedCounts{
		SQLNigori:      0,
		SQLBookmark:    3,
		DynamoNigori:   2,
		DynamoBookmark: 0,
	})

	suite.assertSQLMigrationStatus(bookmarkType, true, true)
	suite.assertSQLMigrationStatus(nigoriType, false, false)
}

func (suite *CommandMigrateTestSuite) TestMigrateDisabled() {
	suite.createSQLDB([]int32{})
	entries := []*sync_pb.SyncEntity{
		getCommitEntity("id1_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id2_bookmark", 0, false, getBookmarkSpecifics()),
		getCommitEntity("id3_nigori", 0, false, getNigoriSpecifics()),
		getCommitEntity("id4_nigori", 0, false, getNigoriSpecifics()),
	}
	msg := getClientToServerCommitMsg(entries)
	rsp := &sync_pb.ClientToServerResponse{}

	// Commit initial entities
	suite.Require().NoError(
		command.HandleClientToServerMessage(suite.cache, msg, rsp, suite.dynamoDB, suite.sqlDB, testClientID),
		"HandleClientToServerMessage should succeed")

	// GetUpdates message
	marker := getMarker(MarkerTokens{
		Nigori:   aws.Int64(math.MaxInt64 - 1000),
		Bookmark: aws.Int64(math.MaxInt64 - 1000),
	})
	msg = getClientToServerGUMsg(marker, sync_pb.SyncEnums_GU_TRIGGER, false, nil)

	// Initial counts
	initialCounts := ExpectedCounts{
		SQLNigori:      0,
		SQLBookmark:    0,
		DynamoNigori:   2,
		DynamoBookmark: 2,
	}
	suite.assertDatastoreCounts(initialCounts)

	// Send multiple GetUpdates messages
	for i := 0; i < 5; i++ {
		suite.sendMessageAndAssertEmptyResponse(msg)

		// Assert that counts haven't changed
		suite.assertDatastoreCounts(initialCounts)
	}

	suite.assertSQLMigrationStatus(bookmarkType, false, false)
	suite.assertSQLMigrationStatus(nigoriType, false, false)
}

// test migration of only one type

func TestCommandMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(CommandMigrateTestSuite))
}
