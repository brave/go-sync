package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/brave-intl/bat-go/middleware"
	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/brave-intl/bat-go/utils/logging"
	"github.com/brave/go-sync/controller"
	"github.com/brave/go-sync/datastore"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func setupLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	return logging.SetupLogger(context.WithValue(ctx, appctx.EnvironmentCTXKey, os.Getenv("ENV")))
}

func setupRouter(ctx context.Context, logger *zerolog.Logger) (context.Context, *chi.Mux) {
	r := chi.NewRouter()

	r.Use(chiware.RequestID)
	r.Use(chiware.RealIP)
	r.Use(chiware.Heartbeat("/"))

	if logger != nil {
		// Also handles panic recovery
		r.Use(hlog.NewHandler(*logger))
		r.Use(hlog.UserAgentHandler("user_agent"))
		r.Use(hlog.RequestIDHandler("req_id", "Request-Id"))
		r.Use(middleware.RequestLogger(logger))
	}

	r.Use(chiware.Timeout(60 * time.Second))
	r.Use(middleware.BearerToken)
	r.Use(middleware.RateLimiter(ctx))

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

	// Setup Sentry.
	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{Dsn: sentryDsn})
		if err != nil {
			logger.Panic().Err(err).Msg("Init sentry failed")
		}
	}

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
