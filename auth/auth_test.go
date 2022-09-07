package auth_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/odedlaz/go-sync/auth"
	"github.com/odedlaz/go-sync/auth/authtest"
	"github.com/odedlaz/go-sync/utils"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func (suite *AuthTestSuite) TestAuthenticate() {
	// invalid token format
	id, err := auth.Authenticate(base64.URLEncoding.EncodeToString([]byte("||")))
	suite.Require().Error(err, "invalid token format should fail")
	suite.Require().Equal("", id, "empty clientID should be returned")

	// invalid signature
	_, tokenHex, _, err := authtest.GenerateToken(utils.UnixMilli(time.Now()))
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(base64.URLEncoding.EncodeToString([]byte("12" + tokenHex)))
	suite.Require().Error(err, "invalid signature should fail")
	suite.Require().Equal("", id)

	// valid token
	tkn, _, expectedID, err := authtest.GenerateToken(utils.UnixMilli(time.Now()))
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(tkn)
	suite.Require().NoError(err, "valid token should succeed")
	suite.Require().Equal(expectedID, id)

	// token expired -1 and +1 day
	tkn, _, _, err = authtest.GenerateToken(utils.UnixMilli(time.Now()) - auth.TokenMaxDuration - 1)
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(tkn)
	suite.Require().Error(err, "outdated token should failed")
	suite.Require().Equal("", id)

	tkn, _, _, err = authtest.GenerateToken(utils.UnixMilli(time.Now()) + auth.TokenMaxDuration + 100)
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(tkn)
	suite.Require().Error(err, "outdated token should failed")
	suite.Require().Equal("", id)
}

func (suite *AuthTestSuite) TestAuthorize() {
	req, err := http.NewRequest("POST", "url", nil)
	suite.Require().NoError(err, "NewRequest should succeed")

	validToken, _, validClientID, err := authtest.GenerateToken(utils.UnixMilli(time.Now()))
	suite.Require().NoError(err, "generate token should succeed")
	outdatedToken, _, _, err := authtest.GenerateToken(utils.UnixMilli(time.Now()) - auth.TokenMaxDuration - 1)
	suite.Require().NoError(err, "generate token should succeed")

	invalidTokenErr := fmt.Errorf("Not a valid token")
	outdatedErr := fmt.Errorf("error authorizing: %w", fmt.Errorf("token is expired"))
	tests := map[string]struct {
		token    string
		clientID string
		err      error
	}{
		"invalid header format": {
			token:    "Bear ",
			clientID: "",
			err:      invalidTokenErr,
		},
		"empty token": {
			token:    "Bearer ",
			clientID: "",
			err:      invalidTokenErr,
		},
		"valid token": {
			token:    "Bearer " + validToken,
			clientID: validClientID,
			err:      nil,
		},
		"outdated token": {
			token:    "Bearer " + outdatedToken,
			clientID: "",
			err:      outdatedErr,
		},
	}
	for testName, test := range tests {
		req.Header.Set("Authorization", test.token)
		clientID, err := auth.Authorize(req)
		suite.Require().Equal(test.err, err,
			"error mismatched for %s test case", testName)
		suite.Require().Equal(test.clientID, clientID,
			"clientID mismatched for %s test case", testName)
	}
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
