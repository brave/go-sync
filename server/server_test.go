package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/odedlaz/go-sync/server"
	"github.com/stretchr/testify/assert"
)

var handler http.Handler

func init() {
	testCtx, logger := server.SetupLogger(context.Background())
	serverCtx, mux := server.SetupRouter(testCtx, logger)
	handler = chi.ServerBaseContext(serverCtx, mux)
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expected := "."
	actual, err := ioutil.ReadAll(rr.Result().Body)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestCommand(t *testing.T) {
	req, err := http.NewRequest("POST", "/v2/command/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.NotEmpty(t, rr.Result().Header.Get("Sane-Time-Millis"))
}
