package server_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brave/go-sync/server"
)

var (
	mux       *chi.Mux
	serverCtx context.Context
)

func init() {
	testCtx, logger := server.SetupLogger(context.Background())
	serverCtx, mux = server.SetupRouter(testCtx, logger)
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req.WithContext(serverCtx))
	assert.Equal(t, http.StatusOK, rr.Code)

	expected := "."
	actual, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestCommand(t *testing.T) {
	req, err := http.NewRequest("POST", "/v2/command/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req.WithContext(serverCtx))
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.NotEmpty(t, rr.Result().Header.Get("Sane-Time-Millis"))
}
