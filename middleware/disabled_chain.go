package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	syncContext "github.com/brave/go-sync/context"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
)

// DisabledChain is a middleware to check for disabled sync chains referenced in a request,
// ending the request early in the case that the request is made against one.
func DisabledChain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID, ok := ctx.Value(syncContext.ContextKeyClientID).(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		db, ok := ctx.Value(syncContext.ContextKeyDatastore).(datastore.Datastore)
		if !ok {
			http.Error(w, "unable to complete request", http.StatusInternalServerError)
			return
		}

		disabled, err := db.IsSyncChainDisabled(clientID)

		if err != nil {
			http.Error(w, "unable to complete request", http.StatusInternalServerError)
			return
		}

		if disabled {
			errCode := sync_pb.SyncEnums_DISABLED_BY_ADMIN
			csRsp := sync_pb.ClientToServerResponse{
				ErrorCode: &errCode,
			}
			out, err := proto.Marshal(&csRsp)
			if err != nil {
				log.Error().Err(err).Msg("Marshall ClientToServerResponse failed")
				http.Error(w, "Marshal Error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(out)
			if err != nil {
				log.Error().Err(err).Msg("Write HTTP response body failed")
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
