package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/brave/go-sync/auth/authtest"
	"github.com/brave/go-sync/datastore/datastoretest"
	"github.com/brave/go-sync/middleware"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	syncContext "github.com/brave/go-sync/synccontext"
)

type MiddlewareTestSuite struct {
	suite.Suite
}

func (suite *MiddlewareTestSuite) TestDisabledChainMiddleware() {
	clientID := "0"

	// Active Chain
	datastore := new(datastoretest.MockDatastore)
	ctx := context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
	datastore.On("IsSyncChainDisabled", ctx, clientID).Return(false, nil)
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler := middleware.DisabledChain(next)
	req, err := http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)

	// Disabled chain
	datastore = new(datastoretest.MockDatastore)
	ctx = context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
	datastore.On("IsSyncChainDisabled", ctx, clientID).Return(true, nil)
	next = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		suite.Fail("Should not reach this point")
	})
	handler = middleware.DisabledChain(next)
	req, err = http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)
	rsp := &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(proto.Unmarshal(rr.Body.Bytes(), rsp), "Unmarshal response should succeed")
	suite.Require().Equal(sync_pb.SyncEnums_DISABLED_BY_ADMIN, rsp.GetErrorCode())

	// DB error
	datastore = new(datastoretest.MockDatastore)
	ctx = context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
	datastore.On("IsSyncChainDisabled", ctx, clientID).Return(false, errors.New("unable to query db"))
	next = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	handler = middleware.DisabledChain(next)
	req, err = http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
	suite.Require().NoError(err, "NewRequestWithContext should succeed")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *MiddlewareTestSuite) TestDisallowNewChainsMiddleware() {
	clientID := "0"
	newRequest := func(isInitialized bool, err error) (*http.Request, *datastoretest.MockDatastore) {
		datastore := new(datastoretest.MockDatastore)
		ctx := context.WithValue(context.Background(), syncContext.ContextKeyClientID, clientID)
		ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, datastore)
		datastore.On("HasServerDefinedUniqueTag", ctx, clientID, "google_chrome_nigori").Return(isInitialized, err)
		req, reqErr := http.NewRequestWithContext(ctx, "POST", "v2/command/", bytes.NewBuffer([]byte{}))
		suite.Require().NoError(reqErr, "NewRequestWithContext should succeed")
		return req, datastore
	}

	// Initialized chain
	req, datastore := newRequest(true, nil)
	nextCalled := false
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		nextCalled = true
	})
	rr := httptest.NewRecorder()
	middleware.DisallowNewChains(next).ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)
	suite.Require().True(nextCalled, "Request for an initialized chain should be passed on")
	datastore.AssertExpectations(suite.T())

	// New chain
	req, datastore = newRequest(false, nil)
	next = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		suite.Fail("Should not reach this point")
	})
	rr = httptest.NewRecorder()
	middleware.DisallowNewChains(next).ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusOK, rr.Code)
	rsp := &sync_pb.ClientToServerResponse{}
	suite.Require().NoError(proto.Unmarshal(rr.Body.Bytes(), rsp), "Unmarshal response should succeed")
	suite.Require().Equal(sync_pb.SyncEnums_DISABLED_BY_ADMIN, rsp.GetErrorCode())
	datastore.AssertExpectations(suite.T())

	// DB error
	req, datastore = newRequest(false, errors.New("unable to query db"))
	rr = httptest.NewRecorder()
	middleware.DisallowNewChains(next).ServeHTTP(rr, req)
	suite.Require().Equal(http.StatusInternalServerError, rr.Code)
	datastore.AssertExpectations(suite.T())
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
