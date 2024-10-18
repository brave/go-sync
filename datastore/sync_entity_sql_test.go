package datastore_test

import (
	"testing"
	"time"

	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SyncEntitySQLTestSuite struct {
	suite.Suite
	sqlDB *datastore.SQLDB
}

func (suite *SyncEntitySQLTestSuite) SetupSuite() {
	var err error
	suite.sqlDB, err = datastore.NewSQLDB(true)
	suite.Require().NoError(err, "Failed to create SQL database")
}

func (suite *SyncEntitySQLTestSuite) SetupTest() {
	err := datastoretest.ResetSQLTables(suite.sqlDB)
	suite.Require().NoError(err, "Failed to reset SQL tables")
}

func createSyncEntity(dataType int32, mtime int64) datastore.SyncEntity {
	id, _ := uuid.NewV7()
	return datastore.SyncEntity{
		ID:        id.String(),
		Version:   &[]int64{1}[0],
		Ctime:     &[]int64{12345678}[0],
		Mtime:     &mtime,
		DataType:  &[]int{int(dataType)}[0],
		Folder:    &[]bool{false}[0],
		Deleted:   &[]bool{false}[0],
		Specifics: []byte{1, 2, 3},
	}
}

func (suite *SyncEntitySQLTestSuite) TestInsertSyncEntity() {
	entity := createSyncEntity(123, 12345678)
	entity.ClientDefinedUniqueTag = &[]string{"test1"}[0]

	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, "client1")
	suite.Require().NoError(err, "GetAndLockChainID should succeed")
	entity.ChainID = chainID

	conflict, err := suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{&entity})
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	suite.Assert().False(conflict, "Insert should not conflict")

	err = tx.Commit()
	suite.Require().NoError(err, "Commit should succeed")

	// Try to insert the same entity again
	tx, err = suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	id, _ := uuid.NewV7()
	entity.ID = id.String()
	conflict, err = suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{&entity})
	suite.Require().NoError(err, "InsertSyncEntity should succeed")
	suite.Assert().True(conflict, "Insert should conflict")

	err = tx.Rollback()
	suite.Require().NoError(err, "Rollback should succeed")
}

func (suite *SyncEntitySQLTestSuite) TestHasItem() {
	entity := createSyncEntity(123, 12345678)
	entity.ClientDefinedUniqueTag = &[]string{"test1"}[0]

	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, "client1")
	suite.Require().NoError(err, "GetAndLockChainID should succeed")
	entity.ChainID = chainID

	_, err = suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{&entity})
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	exists, err := suite.sqlDB.HasItem(tx, *chainID, *entity.ClientDefinedUniqueTag)
	suite.Require().NoError(err, "HasItem should succeed")
	suite.Assert().True(exists, "Item should exist")

	exists, err = suite.sqlDB.HasItem(tx, *chainID, "non_existent_tag")
	suite.Require().NoError(err, "HasItem should succeed")
	suite.Assert().False(exists, "Item should not exist")

	err = tx.Commit()
	suite.Require().NoError(err, "Commit should succeed")
}

func (suite *SyncEntitySQLTestSuite) TestUpdateSyncEntity() {
	entity := createSyncEntity(123, 12345678)
	entity.Specifics = []byte{1, 2}

	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, "client1")
	suite.Require().NoError(err, "GetAndLockChainID should succeed")
	entity.ChainID = chainID

	_, err = suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{&entity})
	suite.Require().NoError(err, "InsertSyncEntity should succeed")

	// Test normal update
	updatedEntity := entity
	updatedEntity.Version = &[]int64{2}[0]
	updatedEntity.Mtime = &[]int64{23456789}[0]
	updatedEntity.Folder = &[]bool{true}[0]

	// Test updating with wrong chain ID
	wrongChainEntity := updatedEntity
	wrongChainEntity.ChainID = &[]int64{*chainID + 1}[0]
	conflict, deleted, err := suite.sqlDB.UpdateSyncEntity(tx, &wrongChainEntity, *entity.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Update should conflict due to wrong chain ID")
	suite.Assert().False(deleted, "Entity should not be deleted")

	// Valid update
	conflict, deleted, err = suite.sqlDB.UpdateSyncEntity(tx, &updatedEntity, *entity.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Update should not conflict")
	suite.Assert().False(deleted, "Entity should not be deleted")

	*entity.Version = *updatedEntity.Version

	*updatedEntity.Version = 3

	// Test updating with wrong version
	conflictEntity := updatedEntity
	conflict, deleted, err = suite.sqlDB.UpdateSyncEntity(tx, &conflictEntity, 99)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Update should conflict due to version mismatch")
	suite.Assert().False(deleted, "Entity should not be deleted")

	// Test updating to deleted state
	*updatedEntity.Deleted = true
	conflict, deleted, err = suite.sqlDB.UpdateSyncEntity(tx, &updatedEntity, *entity.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().False(conflict, "Update should not conflict")
	suite.Assert().True(deleted, "Entity should be deleted")

	*entity.Version = *updatedEntity.Version

	// Test updating a deleted entity
	*updatedEntity.Version = 4
	*updatedEntity.Deleted = false
	conflict, deleted, err = suite.sqlDB.UpdateSyncEntity(tx, &updatedEntity, *entity.Version)
	suite.Require().NoError(err, "UpdateSyncEntity should succeed")
	suite.Assert().True(conflict, "Update should conflict")
	suite.Assert().False(deleted, "Entity should not be deleted")

	err = tx.Commit()
	suite.Require().NoError(err, "Commit should succeed")
}

func (suite *SyncEntitySQLTestSuite) TestGetUpdatesForType() {
	entities := []datastore.SyncEntity{
		createSyncEntity(123, 12345678),
		createSyncEntity(123, 12345679),
		createSyncEntity(123, 12345680),
		createSyncEntity(124, 12345680),
	}

	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	chainID, err := suite.sqlDB.GetAndLockChainID(tx, "client1")
	suite.Require().NoError(err, "GetAndLockChainID should succeed")

	for i := range entities {
		entities[i].ChainID = chainID
		_, err = suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{&entities[i]})
		suite.Require().NoError(err, "InsertSyncEntity should succeed")
	}

	hasChangesRemaining, syncItems, err := suite.sqlDB.GetUpdatesForType(tx, 123, 0, true, *chainID, 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().False(hasChangesRemaining, "Should not have changes remaining")
	suite.Assert().Equal(entities[:3], syncItems)

	hasChangesRemaining, syncItems, err = suite.sqlDB.GetUpdatesForType(tx, 123, 12345678, true, *chainID, 100)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().False(hasChangesRemaining, "Should not have changes remaining")
	suite.Assert().Equal(entities[1:3], syncItems)

	hasChangesRemaining, syncItems, err = suite.sqlDB.GetUpdatesForType(tx, 123, 0, true, *chainID, 2)
	suite.Require().NoError(err, "GetUpdatesForType should succeed")
	suite.Assert().True(hasChangesRemaining, "Should have changes remaining")
	suite.Assert().Equal(entities[:2], syncItems)

	err = tx.Commit()
	suite.Require().NoError(err, "Commit should succeed")
}

func (suite *SyncEntitySQLTestSuite) TestDeleteChain() {
	entity1 := createSyncEntity(123, 12345678)
	entity2 := createSyncEntity(123, 12345678)

	// Insert data for two chains
	tx, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	chainID1, err := suite.sqlDB.GetAndLockChainID(tx, "client1")
	suite.Require().NoError(err, "GetAndLockChainID should succeed for client1")
	entity1.ChainID = chainID1

	chainID2, err := suite.sqlDB.GetAndLockChainID(tx, "client2")
	suite.Require().NoError(err, "GetAndLockChainID should succeed for client2")
	entity2.ChainID = chainID2

	_, err = suite.sqlDB.InsertSyncEntities(tx, []*datastore.SyncEntity{&entity1, &entity2})
	suite.Require().NoError(err, "InsertSyncEntities should succeed")

	err = tx.Commit()
	suite.Require().NoError(err, "Commit should succeed")

	// Delete chain for client1
	tx, err = suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction should succeed")
	defer tx.Rollback()

	err = suite.sqlDB.DeleteChain(tx, *chainID1)
	suite.Require().NoError(err, "DeleteChain should succeed")

	err = tx.Commit()
	suite.Require().NoError(err, "Commit should succeed")

	// Verify that the chain and its entities are deleted for client1
	suite.checkChainExistence(*chainID1, false)

	// Verify that data still exists for client2
	suite.checkChainExistence(*chainID2, true)
}

func (suite *SyncEntitySQLTestSuite) checkChainExistence(chainID int64, shouldExist bool) {
	var expectedCount int
	var count int
	if shouldExist {
		expectedCount = 1
	}
	err := suite.sqlDB.QueryRow("SELECT COUNT(*) FROM entities WHERE chain_id = $1", chainID).Scan(&count)
	suite.Require().NoError(err, "Count query should succeed for entities")
	suite.Assert().Equal(expectedCount, count, "Entities for chain should be correct amount")

	err = suite.sqlDB.QueryRow("SELECT COUNT(*) FROM chains WHERE id = $1", chainID).Scan(&count)
	suite.Require().NoError(err, "Count query should succeed")
	suite.Assert().Equal(expectedCount, count, "Chain entry should be correct amount")
}

func (suite *SyncEntitySQLTestSuite) TestConcurrentGetAndLockChainID() {
	clientID := "testClient"

	// Start first transaction
	tx1, err := suite.sqlDB.Beginx()
	suite.Require().NoError(err, "Begin transaction 1 should succeed")
	defer tx1.Rollback()

	// Get and lock chain ID in first transaction
	chainID1, err := suite.sqlDB.GetAndLockChainID(tx1, clientID)
	suite.Require().NoError(err, "GetAndLockChainID should succeed for tx1")

	// Try to get and lock chain ID in second transaction
	// This should block until the first transaction is committed
	stepChan := make(chan bool)
	go func() {
		// Start second transaction
		tx2, err := suite.sqlDB.Beginx()
		suite.Require().NoError(err, "Begin transaction 2 should succeed")
		defer tx2.Rollback()

		stepChan <- true
		chainID2, err := suite.sqlDB.GetAndLockChainID(tx2, clientID)
		suite.Require().NoError(err, "GetAndLockChainID should succeed for tx2")
		suite.Assert().Equal(*chainID1, *chainID2, "Chain IDs should be the same")

		err = tx2.Commit()
		suite.Require().NoError(err, "Commit transaction 2 should succeed")
		stepChan <- true
	}()

	// Wait until second transaction has started
	<-stepChan

	select {
	case <-stepChan:
		suite.FailNow("Second transaction goroutine exited prematurely")
	case <-time.After(200 * time.Millisecond):
	}

	// Commit the first transaction
	err = tx1.Commit()
	suite.Require().NoError(err, "Commit transaction 1 should succeed")

	// Wait for the second transaction to complete
	select {
	case <-stepChan:
		// Success, second transaction completed
	case <-time.After(5 * time.Second):
		suite.Fail("Second transaction did not complete in time")
	}

}

func TestSyncEntitySQLTestSuite(t *testing.T) {
	suite.Run(t, new(SyncEntitySQLTestSuite))
}
