package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/brave/go-sync/command"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/schema/protobuf/sync_pb"
	syncContext "github.com/brave/go-sync/synccontext"
)

// DisallowNewChains is a middleware to reject requests made against sync
// chains which have not been initialized on this server yet. Existing sync
// chains keep working as usual, including adding new devices to them, but no
// new sync chain can be created.
func DisallowNewChains(next http.Handler) http.Handler {
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

		initialized, err := command.IsSyncChainInitialized(ctx, db, clientID)
		if err != nil {
			log.Error().Err(err).Msg("Check if sync chain is initialized failed")
			http.Error(w, "unable to complete request", http.StatusInternalServerError)
			return
		}

		if !initialized {
			writeErrorCodeResponse(w, sync_pb.SyncEnums_DISABLED_BY_ADMIN)
			return
		}

		next.ServeHTTP(w, r)
	})
}
