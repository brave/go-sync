package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/brave/go-sync/auth"
	syncContext "github.com/brave/go-sync/context"
)

// Auth verifies the token provided is valid, and sets the client id in context
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID, err := auth.Authorize(r)
		if clientID == "" {
			if err != nil {
				log.Error().Err(err).Msg("Authorization failed")
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), syncContext.ContextKeyClientID, clientID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
