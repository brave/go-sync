package middleware

import (
	"context"
	"net/http"

	"github.com/brave/go-sync/auth"
	"github.com/rs/zerolog/log"
)

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

		ctx := context.WithValue(r.Context(), "clientID", clientID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
