package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/brave/go-sync/utils"
)

const (
	saneTimeMillsHeaderKey = "Sane-Time-Millis"
)

// ResponseHeader is a middleware to apply common response headers.
func ResponseHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(saneTimeMillsHeaderKey,
			strconv.FormatInt(utils.UnixMilli(time.Now()), 10))
		next.ServeHTTP(w, r)
	})
}
