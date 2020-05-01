package controller

import (
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/go-chi/chi"
)

// SyncRouter add routers for command and auth endpoint requests.
func SyncRouter(datastore datastore.Datastore) chi.Router {
	r := chi.NewRouter()
	return r
}
