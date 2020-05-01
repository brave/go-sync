package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/brave-experiments/sync-server/controller"
	"github.com/brave-experiments/sync-server/datastore"
	"github.com/brave-intl/bat-go/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func setupLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	var output io.Writer
	if os.Getenv("ENV") != "local" {
		output = os.Stdout
	} else {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	// Use thread-safe, lock-free, non-blocking writer, log msg could be dropped.
	wr := diode.NewWriter(output, 1000, time.Duration(10*time.Millisecond), func(missed int) {
		fmt.Fprintf(os.Stderr, "logger dropped message count: %d", missed)
	})

	// Always print out timestamp.
	log := zerolog.New(wr).With().Timestamp().Logger()

	debug := os.Getenv("DEBUG")
	if debug == "" || debug == "f" || debug == "n" || debug == "0" {
		log = log.Level(zerolog.InfoLevel)
	}

	return log.WithContext(ctx), &log
}

func setupRouter(ctx context.Context, logger *zerolog.Logger) (context.Context, *chi.Mux) {
	r := chi.NewRouter()

	r.Use(chiware.RequestID)
	r.Use(chiware.RealIP)
	r.Use(chiware.Heartbeat("/"))
	r.Use(chiware.Timeout(60 * time.Second))
	r.Use(middleware.BearerToken)

	if logger != nil {
		// Also handles panic recovery
		r.Use(hlog.NewHandler(*logger))
		r.Use(hlog.UserAgentHandler("user_agent"))
		r.Use(hlog.RequestIDHandler("req_id", "Request-Id"))
		r.Use(middleware.RequestLogger(logger))
	}

	db, err := datastore.NewDynamo()
	if err != nil {
		sentry.CaptureException(err)
		log.Panic().Err(err).Msg("Must be able to init datastore to start")
	}

	r.Mount("/v2", controller.SyncRouter(db))
	r.Get("/metrics", middleware.Metrics())

	return ctx, r
}

// StartServer starts the translate proxy server on port 8195
func StartServer() {
	serverCtx, logger := setupLogger(context.Background())
	subLog := logger.Info().Str("prefix", "main")
	subLog.Msg("Starting server")

	serverCtx, r := setupRouter(serverCtx, logger)

	port := ":8295"
	srv := http.Server{Addr: port, Handler: chi.ServerBaseContext(serverCtx, r)}
	err := srv.ListenAndServe()
	if err != nil {
		sentry.CaptureException(err)
		log.Panic().Err(err).Msg("HTTP server start failed!")
	}
}
