package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/brave/go-sync/auth/authtest"
	syncContext "github.com/brave/go-sync/context"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/middleware"
)

type MiddlewareTestSuite struct {
	suite.Suite
}

func (suite *MiddlewareTestSuite) TestDisabledChainMiddleware() {
	clientID := "0"

	// Active Chain
	datastore := new(datastoretest.MockDatastore)
	datastore.On("IsSyncChainDisabled", clientID).Return(false, nil)
	ctx := context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler := middleware.DisabledChain(next)
	req, err := http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)

	// Disabled chain
	datastore = new(datastoretest.MockDatastore)
	datastore.On("IsSyncChainDisabled", clientID).Return(true, nil)
	ctx = context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
	next = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		suite.Fail("Should not reach this point")
	})
	handler = middleware.DisabledChain(next)
	req, err = http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)

	// DB error
	datastore = new(datastoretest.MockDatastore)
	datastore.On("IsSyncChainDisabled", clientID).Return(false, fmt.Errorf("unable to query db"))
	ctx = context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
	next = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler = middleware.DisabledChain(next)
	rr = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware() {
	// Happy path
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID := ctx.Value(syncContext.ContextKeyClientID)
		suite.NotNil(clientID, "Client ID should be set by auth middleware")
	})
	handler := middleware.Auth(next)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	token, _, _, err := authtest.GenerateToken(time.Now().UnixMilli())
	suite.Require().NoError(err, "generate token should succeed")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Invalid bearer token, unauthorized
	next = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler = middleware.Auth(next)
	ctx = context.Background()
	req, err = http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusUnauthorized, rr.Code)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
