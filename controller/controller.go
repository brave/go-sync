package controller

import (
	"net/http"

	"github.com/brave-experiments/sync-server/auth"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-experiments/sync-server/timestamp"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

// SyncRouter add routers for command and auth endpoint requests.
func SyncRouter(datastore datastore.Datastore) chi.Router {
	r := chi.NewRouter()
	r.Post("/auth", Auth(datastore))
	r.Get("/timestamp", Timestamp)
	return r
}

func sendJSONRsp(body []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(body)
	if err != nil {
		log.Error().Err(err).Msg("Write HTTP response body failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Timestamp returns a current timestamp back to sync clients.
func Timestamp(w http.ResponseWriter, r *http.Request) {
	body, err := timestamp.GetTimestamp()
	if err != nil {
		log.Error().Err(err).Msg("Get timestamp failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONRsp(body, w)
}

// Auth handles authentication requests from sync clients.
func Auth(db datastore.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, body, err := auth.Authenticate(r, db)
		if err != nil {
			log.Error().Err(err).Msg("Authenticate failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendJSONRsp(body, w)
	})
}
