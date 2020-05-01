package auth_test

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/brave-experiments/sync-server/auth"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/datastore/datastoretest"
	"github.com/brave-experiments/sync-server/utils"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	dynamo *datastore.Dynamo
}

func (suite *AuthTestSuite) SetupSuite() {
	datastore.Table = "client-entity-token-test-auth"
	var err error
	suite.dynamo, err = datastore.NewDynamo()
	suite.Require().NoError(err, "Failed to get dynamoDB session")
}

func (suite *AuthTestSuite) SetupTest() {
	suite.Require().NoError(
		datastoretest.ResetTable(suite.dynamo), "Failed to reset table")
}

func (suite *AuthTestSuite) TearDownTest() {
	suite.Require().NoError(
		datastoretest.DeleteTable(suite.dynamo), "Failed to delete table")
}

func (suite *AuthTestSuite) TestAuthenticate() {
	defaultTimestampMaxDuration := *auth.TimestampMaxDuration

	// These values are from previous real request.
	encodedClientID := "F33268C948FF4793E20401ABB2B6D1994F9B7CCC5D6C2DDCBE198594BFC32D7C"
	encodedTimestamp := "31353839333331313238383436"
	validSecret := "F715FBC8BBDA7FEC941AE1B29FB9290D64A6001DD102F2D93BAC99151F9E9A18EF35DBEBBE46BB1645375C61A3E9C43D9523811806DEDBB48D32A16F51267B0C"
	invalidSecret := "ABCDEBC8BBDA7FEC941AE1B29FB9290D64A6001DD102F2D93BAC99151F9E9A18EF35DBEBBE46BB1645375C61A3E9C43D9523811806DEDBB48D32A16F51267B0C" // modified from validSecret

	tests := map[string]struct {
		timestampMaxDuration int64
		clientSecret         string
		err                  error
	}{
		"valid signature and timestamp": {
			timestampMaxDuration: math.MaxInt64,
			clientSecret:         validSecret,
			err:                  nil,
		},
		"valid signature and outdated timestamp": {
			timestampMaxDuration: 0,
			clientSecret:         validSecret,
			err:                  fmt.Errorf("timestamp is outdated"),
		},
		"invalid signature": {
			timestampMaxDuration: math.MaxInt64,
			clientSecret:         invalidSecret,
			err:                  fmt.Errorf("signature verification failed"),
		},
	}

	for testName, test := range tests {
		form := url.Values{
			"client_id":     {encodedClientID},
			"timestamp":     {encodedTimestamp},
			"client_secret": {test.clientSecret},
		}
		req, err := http.NewRequest("POST", "url", strings.NewReader(form.Encode()))
		suite.Require().NoError(err, "NewRequest should succeed")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		*auth.TimestampMaxDuration = test.timestampMaxDuration
		token, rsp, err := auth.Authenticate(req, suite.dynamo)
		suite.Assert().Equal(test.err, err, "err mismatch for %s test case", testName)

		if test.err == nil {
			suite.Assert().NotEqual("", token, "%s: success request should not return empty token", testName)

			authRsp := auth.Response{AccessToken: token, ExpiresIn: auth.TokenMaxDuration}
			outRsp, err := json.Marshal(authRsp)
			suite.Require().NoError(err, "json marshal should succeed")
			suite.Assert().Equal(outRsp, rsp, "rsp mismatch for %s test case", testName)
		} else {
			suite.Assert().Equal("", token, "%s: fail request should return empty token", testName)
			suite.Assert().Nil(rsp, "response should be nil for %s test case", testName)
		}
	}

	*auth.TimestampMaxDuration = defaultTimestampMaxDuration
}

func (suite *AuthTestSuite) TestAuthorize() {
	outdatedTime := utils.UnixMilli(time.Unix(0, 0))
	validTime := utils.UnixMilli(time.Now().Add(time.Minute * 30))
	suite.Require().NoError(suite.dynamo.InsertClientToken("key", "token1", outdatedTime))
	suite.Require().NoError(suite.dynamo.InsertClientToken("key", "token2", validTime))

	req, err := http.NewRequest("POST", "url", nil)
	suite.Require().NoError(err, "NewRequest should succeed")

	invalidTokenErr := fmt.Errorf("Not a valid token")
	tests := map[string]struct {
		header   string
		clientID string
		err      error
	}{
		"invalid header format": {
			header:   "Bear ",
			clientID: "",
			err:      invalidTokenErr,
		},
		"empty token": {
			header:   "Bearer ",
			clientID: "",
			err:      invalidTokenErr,
		},
		"valid token": {
			header:   "Bearer token2",
			clientID: "key",
			err:      nil,
		},
		"outdated token": {
			header:   "Bearer token1",
			clientID: "",
			err:      invalidTokenErr,
		},
		"invalid token": {
			header:   "Bearer token3",
			clientID: "",
			err:      invalidTokenErr,
		},
	}

	for testName, test := range tests {
		req.Header.Set("Authorization", test.header)
		clientID, err := auth.Authorize(suite.dynamo, req)
		suite.Require().Equal(test.err, err,
			"error mismatched for %s test case", testName)
		suite.Require().Equal(test.clientID, clientID,
			"clientID mismatched for %s test case", testName)
	}
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
