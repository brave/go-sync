package controller

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/brave-intl/bat-go/middleware"
	"github.com/brave-intl/bat-go/utils/closers"
	"github.com/brave/go-sync/auth"
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
func SyncRouter(datastore datastore.Datastore) chi.Router {
	r := chi.NewRouter()
	r.Method("POST", "/command/", middleware.InstrumentHandler("Command", Command(datastore)))
	return r
}

// Command handles GetUpdates and Commit requests from sync clients.
func Command(db datastore.Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorize
		clientID, err := auth.Authorize(r)
		if clientID == "" {
			if err != nil {
				log.Error().Err(err).Msg("Authorization failed")
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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
		err = command.HandleClientToServerMessage(pb, pbRsp, db, clientID)
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
