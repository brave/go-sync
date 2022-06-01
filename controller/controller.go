package controller

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/brave-intl/bat-go/middleware"
	"github.com/brave-intl/bat-go/utils/closers"
	"github.com/brave/go-sync/account"
	"github.com/brave/go-sync/auth"
	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

const (
	payloadLimit10MB = 1024 * 1024 * 10
)

// SyncRouter add routers for command and auth endpoint requests.
func SyncRouter(cache *cache.Cache, datastore datastore.Datastore) chi.Router {
	r := chi.NewRouter()
	r.Method("POST", "/command/", middleware.InstrumentHandler("Command", Command(cache, datastore)))
	r.Method("DELETE", "/account/{clientId}", middleware.InstrumentHandler("Account", Account(cache, datastore)))
	return r
}

// Account handles delete account requests from sync clients.
func Account(cache *cache.Cache, db datastore.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID, ok := ctx.Value("clientID").(string)
		if !ok {
			http.Error(w, "Missing client id", http.StatusUnauthorized)
			return
		}

		err := account.HandleDelete(clientID, db)
		if err != nil {
			http.Error(w, "Error deleting account", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func authorize(w http.ResponseWriter, r *http.Request) (string, error) {
	clientID, err := auth.Authorize(r)
	if clientID == "" {
		if err != nil {
			log.Error().Err(err).Msg("Authorization failed")
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return "", err
	}
	return clientID, nil
}

// Command handles GetUpdates and Commit requests from sync clients.
func Command(cache *cache.Cache, db datastore.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID, ok := ctx.Value("clientID").(string)
		if !ok {
			http.Error(w, "missing client id", http.StatusUnauthorized)
			return
		}

		reader := r.Body
		// Create a gzip reader if needed.
		if r.Header.Get("Content-Encoding") == "gzip" {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Error().Err(err).Msg("Create gzip reader failed")
				http.Error(w, "Create gzip reader failed", http.StatusInternalServerError)
				return
			}
			defer closers.Panic(gr)
			reader = gr
		}

		msg, err := ioutil.ReadAll(io.LimitReader(reader, payloadLimit10MB))
		if err != nil {
			log.Error().Err(err).Msg("Read request body failed")
			http.Error(w, "Read request body error", http.StatusInternalServerError)
			return
		}

		// Unmarshal into ClientToServerMessage
		pb := &sync_pb.ClientToServerMessage{}
		err = proto.Unmarshal(msg, pb)
		if err != nil {
			log.Error().Err(err).Msg("Unmarshall ClientToServerMessage failed")
			http.Error(w, "Unmarshal error", http.StatusInternalServerError)
			return
		}

		pbRsp := &sync_pb.ClientToServerResponse{}
		err = command.HandleClientToServerMessage(cache, pb, pbRsp, db, clientID)
		if err != nil {
			log.Error().Err(err).Msg("Handle command message failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		out, err := proto.Marshal(pbRsp)
		if err != nil {
			log.Error().Err(err).Msg("Marshall ClientToServerResponse failed")
			http.Error(w, "Marshal Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(out)
		if err != nil {
			log.Error().Err(err).Msg("Write HTTP response body failed")
			return
		}
	})
}
