package server_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brave-experiments/sync-server/server"
	"github.com/go-chi/chi"
)

var handler http.Handler

func init() {
	testCtx, logger := server.SetupLogger(context.Background())
	serverCtx, mux := server.SetupRouter(testCtx, logger)
	handler = chi.ServerBaseContext(serverCtx, mux)
}

func TestPing(t *testing.T) {
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	expected := "."
	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(actual) {
		t.Errorf("Expected the message '%s'\n", expected)
	}
}
