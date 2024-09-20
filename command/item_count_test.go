package command_test

import (
	"context"
	"testing"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type ItemCountTestSuite struct {
	suite.Suite
	dynamoDB *datastore.Dynamo
	cache    *cache.Cache
	sqlDB    *datastore.SQLDB
}

func (suite *ItemCountTestSuite) SetupSuite() {
	var rollouts string
	suite.T().Setenv(datastore.SQLSaveRolloutsEnvKey, rollouts)
	suite.T().Setenv(datastore.SQLMigrateRolloutsEnvKey, rollouts)

	datastore.Table = "client-entity-test-command"
	var err error
	suite.dynamoDB, err = datastore.NewDynamo(true)
	suite.Require().NoError(err, "Failed to get dynamoDB session")
	suite.sqlDB, err = datastore.NewSQLDB(true)
	suite.Require().NoError(err, "Failed to get SQL DB session")

	suite.cache = cache.NewCache(cache.NewRedisClient())
}

func (suite *ItemCountTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetDynamoTable(suite.dynamoDB), "Failed to reset Dynamo table")
	suite.Require().NoError(
		datastoretest.ResetSQLTables(suite.sqlDB), "Failed to reset SQL tables")
}

func (suite *ItemCountTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamoDB), "Failed to delete table")
	suite.Require().NoError(
		suite.cache.FlushAll(context.Background()), "Failed to clear cache")
}

func (suite *ItemCountTestSuite) insertSyncEntity(tx *sqlx.Tx, itemCounts *command.ItemCounts, insertInSQL bool, dataType int, clientID string, chainID int64) *datastore.SyncEntity {
	id, err := uuid.NewV7()
	suite.Require().NoError(err, "Failed to generate UUID")

	entity := &datastore.SyncEntity{
		ChainID:                &chainID,
		ClientID:               clientID,
		ID:                     id.String(),
		DataType:               &dataType,
		Version:                &[]int64{1}[0],
		Mtime:                  &[]int64{123}[0],
		Ctime:                  &[]int64{123}[0],
		Specifics:              []byte{1, 2},
		Folder:                 &[]bool{false}[0],
		Deleted:                &[]bool{false}[0],
		ClientDefinedUniqueTag: &[]string{id.String()}[0],
		DataTypeMtime:          &[]string{"123#12345678"}[0],
	}

	if insertInSQL {
		conflict, err := suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{entity})
		suite.Require().NoError(err, "Failed to insert sync entity in SQL")
		suite.Require().False(conflict, "Unexpected conflict when inserting sync entity in SQL")
	} else {
		conflict, err := suite.dynamoDB.InsertSyncEntity(entity)
		suite.Require().NoError(err, "Failed to insert sync entity in DynamoDB")
		suite.Require().False(conflict, "Unexpected conflict when inserting sync entity in DynamoDB")
	}
	suite.Require().NoError(itemCounts.RecordChange(dataType, false, insertInSQL), "Should be able record change")
	return entity
}

func (suite *ItemCountTestSuite) deleteSyncEntity(tx *sqlx.Tx, itemCounts *command.ItemCounts, deleteInSQL bool, entity *datastore.SyncEntity) {
	*entity.Version = 2
	*entity.Deleted = true
	if deleteInSQL {
		conflict, deleted, err := suite.sqlDB.UpdateSyncEntity(tx, entity, 1)
		suite.Require().NoError(err, "Failed to delete sync entity in SQL")
		suite.Require().False(conflict, "Unexpected conflict when deleting sync entity in SQL")
		suite.Require().True(deleted, "Expected entity to be marked as deleted in SQL")
	} else {
		conflict, deleted, err := suite.dynamoDB.UpdateSyncEntity(entity, 1)
		suite.Require().NoError(err, "Failed to delete sync entity in DynamoDB")
		suite.Require().False(conflict, "Unexpected conflict when deleting sync entity in DynamoDB")
		suite.Require().True(deleted, "Expected entity to be marked as deleted in DynamoDB")
	}
	suite.Require().NoError(itemCounts.RecordChange(*entity.DataType, true, deleteInSQL), "Should be able to record change")
}

func (suite *ItemCountTestSuite) TestPreloaded() {
	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, "client1")
	suite.Require().NoError(err, "Failed to get chain ID")

	itemCounts, err := command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err)

	suite.Equal(0, itemCounts.SumCounts(false), "Expected initial sum of item counts to be zero")
	suite.Equal(0, itemCounts.SumCounts(true), "Expected initial sum of item counts to be zero")
}

func (suite *ItemCountTestSuite) TestInsertAndCountItems() {
	clientID := "client1"

	// Start a new transaction for insertions
	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction for insertions")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, clientID)
	suite.Require().NoError(err, "Failed to get chain ID")

	itemCounts, err := command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)

	// Insert items
	suite.insertSyncEntity(tx, itemCounts, true, 123, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, true, 124, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, true, datastore.HistoryTypeID, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, true, datastore.HistoryDeleteDirectiveTypeID, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, false, 123, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, false, datastore.HistoryTypeID, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, false, datastore.HistoryDeleteDirectiveTypeID, clientID, *chainID)

	suite.Equal(4, itemCounts.SumCounts(true), "Expected history total count of 4")
	suite.Equal(7, itemCounts.SumCounts(false), "Expected total count of 7")

	suite.Require().NoError(tx.Commit(), "Failed to commit transaction")

	// Start a new transaction for counting
	tx, err = suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction for counting")
	defer tx.Rollback()

	itemCounts, err = command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err, "Failed to get item counts")

	suite.Equal(4, itemCounts.SumCounts(true), "Expected history total count of 4")
	suite.Equal(7, itemCounts.SumCounts(false), "Expected total count of 7")

	clientID = "client2"
	chainID, err = suite.sqlDB.GetAndLockChainID(tx, clientID)
	suite.Require().NoError(err, "Failed to get chain ID for other client")

	otherItemCounts, err := command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err, "Failed to get item counts for other client")

	suite.Equal(0, otherItemCounts.SumCounts(true), "Expected history total count of 0 for other client")
	suite.Equal(0, otherItemCounts.SumCounts(false), "Expected total count of 0 for other client")
}

func (suite *ItemCountTestSuite) TestDeleteAfterInsertCommit() {
	clientID := "client1"

	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, clientID)
	suite.Require().NoError(err, "Failed to get chain ID")

	itemCounts, err := command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err, "Failed to get initial item counts")

	var sqlEntitiesToDelete []*datastore.SyncEntity
	var dynamoEntitiesToDelete []*datastore.SyncEntity

	sqlEntitiesToDelete = append(sqlEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, true, 123, clientID, *chainID))
	suite.insertSyncEntity(tx, itemCounts, true, 124, clientID, *chainID)
	sqlEntitiesToDelete = append(sqlEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, true, datastore.HistoryTypeID, clientID, *chainID))
	suite.insertSyncEntity(tx, itemCounts, true, datastore.HistoryTypeID, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, false, 125, clientID, *chainID)
	dynamoEntitiesToDelete = append(dynamoEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, false, 125, clientID, *chainID))
	dynamoEntitiesToDelete = append(dynamoEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, false, datastore.HistoryDeleteDirectiveTypeID, clientID, *chainID))
	suite.insertSyncEntity(tx, itemCounts, false, datastore.HistoryTypeID, clientID, *chainID)

	suite.Require().NoError(tx.Commit(), "Failed to commit transaction")

	tx, err = suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction")
	defer tx.Rollback()

	for _, entity := range sqlEntitiesToDelete {
		suite.deleteSyncEntity(tx, itemCounts, true, entity)
	}
	for _, entity := range dynamoEntitiesToDelete {
		suite.deleteSyncEntity(tx, itemCounts, false, entity)
	}

	suite.Equal(2, itemCounts.SumCounts(true), "Expected history count of 2 after deletions")
	suite.Equal(4, itemCounts.SumCounts(false), "Expected total count of 4 after deletions")

	suite.Require().NoError(tx.Commit(), "Failed to commit transaction")

	tx, err = suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction for final count")
	defer tx.Rollback()

	itemCounts, err = command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err, "Failed to get final item counts")

	suite.Equal(2, itemCounts.SumCounts(true), "Expected history count of 2 after deletions")
	suite.Equal(4, itemCounts.SumCounts(false), "Expected total count of 4 after deletions")
}

func (suite *ItemCountTestSuite) TestDeleteBeforeInsertCommit() {
	clientID := "client1"

	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, clientID)
	suite.Require().NoError(err, "Failed to get chain ID")

	itemCounts, err := command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err, "Failed to get initial item counts")

	var sqlEntitiesToDelete []*datastore.SyncEntity
	var dynamoEntitiesToDelete []*datastore.SyncEntity

	sqlEntitiesToDelete = append(sqlEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, true, 123, clientID, *chainID))
	suite.insertSyncEntity(tx, itemCounts, true, 124, clientID, *chainID)
	sqlEntitiesToDelete = append(sqlEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, true, datastore.HistoryTypeID, clientID, *chainID))
	suite.insertSyncEntity(tx, itemCounts, true, datastore.HistoryTypeID, clientID, *chainID)
	suite.insertSyncEntity(tx, itemCounts, false, 125, clientID, *chainID)
	dynamoEntitiesToDelete = append(dynamoEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, false, 125, clientID, *chainID))
	dynamoEntitiesToDelete = append(dynamoEntitiesToDelete, suite.insertSyncEntity(tx, itemCounts, false, datastore.HistoryDeleteDirectiveTypeID, clientID, *chainID))
	suite.insertSyncEntity(tx, itemCounts, false, datastore.HistoryTypeID, clientID, *chainID)

	for _, entity := range sqlEntitiesToDelete {
		suite.deleteSyncEntity(tx, itemCounts, true, entity)
	}
	for _, entity := range dynamoEntitiesToDelete {
		suite.deleteSyncEntity(tx, itemCounts, false, entity)
	}

	// Check counts before commit
	suite.Equal(2, itemCounts.SumCounts(true), "Expected SQL count of 2 before commit")
	suite.Equal(4, itemCounts.SumCounts(false), "Expected total count of 4 before commit")

	suite.Require().NoError(tx.Commit(), "Failed to commit transaction")

	// Start a new transaction for final count
	tx, err = suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Failed to start transaction for final count")
	defer tx.Rollback()

	itemCounts, err = command.GetItemCounts(suite.cache, suite.dynamoDB, suite.sqlDB, tx, clientID, *chainID)
	suite.Require().NoError(err, "Failed to get final item counts")

	// Check counts after commit
	suite.Equal(2, itemCounts.SumCounts(true), "Expected SQL count of 2 after commit")
	suite.Equal(4, itemCounts.SumCounts(false), "Expected total count of 4 after commit")
}

func TestItemCountTestSuite(t *testing.T) {
	suite.Run(t, new(ItemCountTestSuite))
}
