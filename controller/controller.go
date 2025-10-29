package controller

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/brave-intl/bat-go/libs/closers"
	"github.com/brave-intl/bat-go/libs/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/brave/go-sync/cache"
	"github.com/brave/go-sync/command"
	syncContext "github.com/brave/go-sync/context"
	"github.com/brave/go-sync/datastore"
	syncMiddleware "github.com/brave/go-sync/middleware"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

const (
	payloadLimit10MB = 1024 * 1024 * 10
)

// SyncRouter add routers for command and auth endpoint requests.
func SyncRouter(cache *cache.Cache, datastore datastore.Datastore) chi.Router {
	r := chi.NewRouter()
	r.Use(syncMiddleware.Auth)
	r.Use(syncMiddleware.DisabledChain)
	r.Method("POST", "/command/", middleware.InstrumentHandler("Command", Command(cache, datastore)))
	return r
}

// Command handles GetUpdates and Commit requests from sync clients.
func Command(cache *cache.Cache, db datastore.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID, ok := ctx.Value(syncContext.ContextKeyClientID).(string)
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
			defer closers.Panic(ctx, gr)
			reader = gr
		}

		msg, err := io.ReadAll(io.LimitReader(reader, payloadLimit10MB))
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
			return
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
