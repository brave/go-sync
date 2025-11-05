package datastore_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
)

type ItemCountTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *ItemCountTestSuite) SetupSuite() {
	datastore.Table = "client-entity-test-datastore"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *ItemCountTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *ItemCountTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
}

func (suite *ItemCountTestSuite) TestGetClientItemCount() {
	// Insert two items for test.
	items := []datastore.ClientItemCounts{
		{ClientID: "client1", ID: "client1", ItemCount: 5},
		{ClientID: "client2", ID: "client2", ItemCount: 10},
	}
	for _, item := range items {
		existing := datastore.ClientItemCounts{ClientID: item.ClientID, ID: item.ID, Version: datastore.CurrentCountVersion}
		suite.Require().NoError(
			suite.dynamo.UpdateClientItemCount(&existing, item.ItemCount, 0))
	}

	for _, item := range items {
		count, err := suite.dynamo.GetClientItemCount(item.ClientID)
		suite.Require().NoError(err, "GetClientItemCount should succeed")
		suite.Equal(count.ItemCount, item.ItemCount, "ItemCount should match")
	}

	// Non-exist client item count should succeed with count = 0.
	count, err := suite.dynamo.GetClientItemCount("client3")
	suite.Require().NoError(err, "Get non-exist ClientItemCount should succeed")
	suite.Equal(0, count.ItemCount)
}

func (suite *ItemCountTestSuite) TestUpdateClientItemCount() {
	items := []datastore.ClientItemCounts{
		{ClientID: "client1", ID: "client1", ItemCount: 1},
		{ClientID: "client1", ID: "client1", ItemCount: 5},
		{ClientID: "client2", ID: "client2", ItemCount: 10},
	}
	expectedItems := []datastore.ClientItemCounts{
		{ClientID: "client1", ID: "client1", ItemCount: 6},
		{ClientID: "client2", ID: "client2", ItemCount: 10},
	}

	for _, item := range items {
		count, err := suite.dynamo.GetClientItemCount(item.ClientID)
		suite.Require().NoError(err)
		suite.Require().NoError(
			suite.dynamo.UpdateClientItemCount(count, item.ItemCount, 0))
	}

	clientCountItems, err := datastoretest.ScanClientItemCounts(suite.dynamo)
	suite.Require().NoError(err, "ScanClientItemCounts should succeed")
	sort.Sort(datastore.ClientItemCountByClientID(clientCountItems))
	sort.Sort(datastore.ClientItemCountByClientID(expectedItems))
	for i := range clientCountItems {
		clientCountItems[i].Version = 0
		clientCountItems[i].LastPeriodChangeTime = 0
	}
	suite.Equal(expectedItems, clientCountItems)
}

func TestItemCountTestSuite(t *testing.T) {
	suite.Run(t, new(ItemCountTestSuite))
}
