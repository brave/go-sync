package main

import (
	"os"

	"github.com/brave/go-sync/server"
)

func main() {
	server.StartServer(server.ServerOpts{SentryDSN: os.Getenv("SENTRY_DSN")})
}
