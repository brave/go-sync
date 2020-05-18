package datastore_test

import (
	"sort"
	"testing"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/datastore/datastoretest"
	"github.com/stretchr/testify/suite"
)

type ItemCountTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *ItemCountTestSuite) SetupSuite() {
	datastore.Table = "client-entity-token-test-datastore"
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
	items := []datastore.ClientItemCount{
		{ClientID: "client1", ID: "client1", ItemCount: 5},
		{ClientID: "client2", ID: "client2", ItemCount: 10},
	}
	for _, item := range items {
		suite.Require().NoError(
			suite.dynamo.UpdateClientItemCount(item.ClientID, item.ItemCount))
	}

	for _, item := range items {
		count, err := suite.dynamo.GetClientItemCount(item.ClientID)
		suite.Require().NoError(err, "GetClientItemCount should succeed")
		suite.Assert().Equal(count, item.ItemCount, "ItemCount should match")
	}

	// Non-exist client item count should succeed with count = 0.
	count, err := suite.dynamo.GetClientItemCount("client3")
	suite.Require().NoError(err, "Get non-exist ClientItemCount should succeed")
	suite.Assert().Equal(count, 0)
}

func (suite *ItemCountTestSuite) TestUpdateClientItemCount() {
	items := []datastore.ClientItemCount{
		{ClientID: "client1", ID: "client1", ItemCount: 1},
		{ClientID: "client1", ID: "client1", ItemCount: 5},
		{ClientID: "client2", ID: "client2", ItemCount: 10},
	}
	expectedItems := []datastore.ClientItemCount{
		{ClientID: "client1", ID: "client1", ItemCount: 6},
		{ClientID: "client2", ID: "client2", ItemCount: 10},
	}

	for _, item := range items {
		suite.Require().NoError(
			suite.dynamo.UpdateClientItemCount(item.ClientID, item.ItemCount))
	}

	clientCountItems, err := datastoretest.ScanClientItemCounts(suite.dynamo)
	suite.Require().NoError(err, "ScanClientItemCounts should succeed")
	sort.Sort(datastore.ClientItemCountByClientID(clientCountItems))
	sort.Sort(datastore.ClientItemCountByClientID(expectedItems))
	suite.Assert().Equal(expectedItems, clientCountItems)
}

func TestItemCountTestSuite(t *testing.T) {
	suite.Run(t, new(ItemCountTestSuite))
}
