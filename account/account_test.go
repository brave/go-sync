package account_test

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/brave/go-sync/account"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/stretchr/testify/suite"
)

const (
	clientID    string = "client"
	deletedType int32  = 3
	pk          string = "ClientID"
	sk          string = "ID"
	projPk      string = "ClientID, ID"
)

type AccountTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *AccountTestSuite) SetupSuite() {
	datastore.Table = "client-entity-test-account"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *AccountTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *AccountTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
}

func (suite *AccountTestSuite) TestHandleDelete() {
	// Write a bunch of values to dynamo
	for i := 0; i < 10; i++ {
		entity := datastore.SyncEntity{
			ClientID: clientID,
			ID:       strconv.FormatInt(int64(i), 10),
		}
		suite.dynamo.InsertSyncEntity(&entity)
	}

	// Check that delete succeeds
	suite.Require().NoError(
		account.HandleDelete(clientID, suite.dynamo),
		"HandleClientDelete should succeed")

	// check that the only remaining value is the disabled marker
	pkb := expression.Key(pk)
	pkv := expression.Value(clientID)
	keyCond := expression.KeyEqual(pkb, pkv)
	exprs := expression.NewBuilder().WithKeyCondition(keyCond)
	expr, err := exprs.Build()
	suite.Require().NoError(err, "Failed to build expression to get updates")

	input := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      aws.String(projPk),
		TableName:                 aws.String(datastore.Table),
	}

	out, err := suite.dynamo.Query(input)
	suite.Require().NoError(err, "Failed to query dynamo")
	count := *out.Count

	suite.Assert().Equal(int64(1), count)
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
