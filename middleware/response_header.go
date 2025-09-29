package middleware

import (
	"net/http"
	"strconv"
	"time"
)

func saneTimeMillsHeaderFunc(w http.ResponseWriter) {
	w.Header().Set("Sane-Time-Millis", strconv.FormatInt(time.Now().UnixMilli(), 10))
}

var (
	headerFuncs = []func(w http.ResponseWriter){
		saneTimeMillsHeaderFunc,
	}
)

// CommonResponseHeaders is a middleware to apply common response headers.
func CommonResponseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, headerFunc := range headerFuncs {
			headerFunc(w)
		}

		next.ServeHTTP(w, r)
	})
}
