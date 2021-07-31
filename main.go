package main

import (
	"os"

	"github.com/brave/go-sync/server"
)

func main() {
	server.StartServer(
		server.Opts{
			SentryDSN: os.Getenv("SENTRY_DSN"),
			Addr:      ":8295",
		},
	)
}
