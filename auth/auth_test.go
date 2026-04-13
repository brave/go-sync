package auth_test

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/brave/go-sync/auth"
	"github.com/brave/go-sync/auth/authtest"
)

type AuthTestSuite struct {
	suite.Suite
}

func (suite *AuthTestSuite) TestBlockedClientIDCaseInsensitive() {
	// Generate a valid token and grab its hex public key.
	_, tokenHex, publicKeyHex, err := authtest.GenerateToken(time.Now().UnixMilli())
	suite.Require().NoError(err)

	// Block the key using uppercase hex in the env var.
	suite.T().Setenv("BLOCKED_CLIENT_IDS", strings.ToUpper(publicKeyHex))

	// Re-encode the token with the public key hex uppercased so that a naive
	// case-sensitive comparison would not match the lowercase blocked entry.
	parts := strings.Split(tokenHex, "|")
	suite.Require().Len(parts, 3)
	parts[2] = strings.ToUpper(parts[2])
	mixedToken := base64.URLEncoding.EncodeToString([]byte(strings.Join(parts, "|")))

	id, err := auth.Authenticate(mixedToken)
	suite.Require().Error(err, "blocked client should be rejected regardless of hex casing")
	suite.Require().Empty(id)
}

func (suite *AuthTestSuite) TestAuthenticate() {
	// invalid token format
	id, err := auth.Authenticate(base64.URLEncoding.EncodeToString([]byte("||")))
	suite.Require().Error(err, "invalid token format should fail")
	suite.Require().Empty(id, "empty clientID should be returned")

	// invalid signature
	_, tokenHex, _, err := authtest.GenerateToken(time.Now().UnixMilli())
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(base64.URLEncoding.EncodeToString([]byte("12" + tokenHex)))
	suite.Require().Error(err, "invalid signature should fail")
	suite.Require().Empty(id)

	// valid token
	tkn, _, expectedID, err := authtest.GenerateToken(time.Now().UnixMilli())
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(tkn)
	suite.Require().NoError(err, "valid token should succeed")
	suite.Require().Equal(expectedID, id)

	// token expired -1 and +1 day
	tkn, _, _, err = authtest.GenerateToken(time.Now().UnixMilli() - auth.TokenMaxDuration - 1)
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(tkn)
	suite.Require().Error(err, "outdated token should failed")
	suite.Require().Empty(id)

	tkn, _, _, err = authtest.GenerateToken(time.Now().UnixMilli() + auth.TokenMaxDuration + 100)
	suite.Require().NoError(err, "generate token should succeed")
	id, err = auth.Authenticate(tkn)
	suite.Require().Error(err, "outdated token should failed")
	suite.Require().Empty(id)
}

func (suite *AuthTestSuite) TestAuthorize() {
	req, err := http.NewRequest("POST", "url", nil)
	suite.Require().NoError(err, "NewRequest should succeed")

	validToken, _, validClientID, err := authtest.GenerateToken(time.Now().UnixMilli())
	suite.Require().NoError(err, "generate token should succeed")
	outdatedToken, _, _, err := authtest.GenerateToken(time.Now().UnixMilli() - auth.TokenMaxDuration - 1)
	suite.Require().NoError(err, "generate token should succeed")

	invalidTokenErr := errors.New("not a valid token")
	outdatedErr := fmt.Errorf("error authorizing: %w", errors.New("token is expired"))
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
