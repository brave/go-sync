package datastore_test

import (
	"sort"
	"testing"
	"time"

	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/datastore/datastoretest"
	"github.com/brave-experiments/sync-server/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *TokenTestSuite) SetupSuite() {
	datastore.Table = "client-entity-token-test-datastore"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *TokenTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *TokenTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
}

func (suite *TokenTestSuite) TestInsertClientToken() {
	clientID := uuid.NewV4().String()
	token := uuid.NewV4().String()
	expireAt := utils.UnixMilli(time.Now())

	// Insert a client should succeed.
	err := suite.dynamo.InsertClientToken(clientID, token, expireAt)
	suite.Require().NoError(err, "Insert a client-token should succeed")

	// Insert the same ClientID, Token should fail.
	err = suite.dynamo.InsertClientToken(clientID, token, expireAt)
	suite.Require().Error(err, "Insert the same client-token should fail")

	// Insert another token for this client should succeed.
	otherToken := uuid.NewV4().String()
	err = suite.dynamo.InsertClientToken(clientID, otherToken, expireAt)
	suite.Require().NoError(err, "Insert a second token for a client should succeed")

	// Check that we have two ClientToken in the DB.
	expectClientTokens := []datastore.ClientToken{
		{
			ClientID: clientID,
			Token:    token,
			ExpireAt: expireAt,
		},
		{
			ClientID: clientID,
			Token:    otherToken,
			ExpireAt: expireAt,
		},
	}
	sort.Slice(expectClientTokens, func(i, j int) bool { // Sort by the range key.
		return expectClientTokens[i].Token < expectClientTokens[j].Token
	})
	clientTokens, err := datastoretest.ScanClientTokens(suite.dynamo)
	suite.Require().NoError(err, "Scan client tokens should succeed")
	suite.Assert().Equal(clientTokens, expectClientTokens)
}

func (suite *TokenTestSuite) TestGetClientID() {
	expiredToken := datastore.ClientToken{
		ClientID: "id1",
		Token:    "expired_token",
		ExpireAt: utils.UnixMilli(time.Now()),
	}
	validToken := datastore.ClientToken{
		ClientID: "id1",
		Token:    "valid_token1",
		ExpireAt: utils.UnixMilli(time.Now()) + 60000,
	}
	otherClientToken := datastore.ClientToken{
		ClientID: "id2",
		Token:    "valid_token2",
		ExpireAt: utils.UnixMilli(time.Now()) + 60000,
	}
	clientTokens := []datastore.ClientToken{expiredToken, validToken, otherClientToken}

	for _, token := range clientTokens {
		err := suite.dynamo.InsertClientToken(token.ClientID, token.Token, token.ExpireAt)
		suite.Require().NoError(err, "Insert a client-token should succeed")
	}

	clientID, err := suite.dynamo.GetClientID("invalid_token")
	suite.Require().Error(err, "Get clientID from an invalid token should fail")
	suite.Assert().Equal(clientID, "")

	clientID, err = suite.dynamo.GetClientID(expiredToken.Token)
	suite.Require().Error(err, "Get clientID from an expired token should fail")
	suite.Assert().Equal(clientID, "")

	clientID, err = suite.dynamo.GetClientID(validToken.Token)
	suite.Require().NoError(err, "Get clientID from a valid token should succeed")
	suite.Assert().Equal(clientID, validToken.ClientID)

	clientID, err = suite.dynamo.GetClientID(otherClientToken.Token)
	suite.Require().NoError(err, "Get clientID from a valid token should succeed")
	suite.Assert().Equal(clientID, otherClientToken.ClientID)
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
