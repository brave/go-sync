package middleware

import (
	"net/http"

	"github.com/brave/go-sync/datastore"
)

// DisabledChain is a middleware to check for disabled sync chains referenced in a request,
// ending the request early in the case that the request is made against one.
func DisabledChain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientID, ok := ctx.Value("clientID").(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		db, ok := ctx.Value("datastore").(datastore.Datastore)
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
			http.Error(w, "sync chain not found", http.StatusGone)
			return
		}

		next.ServeHTTP(w, r)
	})
}
