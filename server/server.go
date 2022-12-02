package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // pprof magic
	"os"
	"time"

	batware "github.com/brave-intl/bat-go/middleware"
	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/brave-intl/bat-go/utils/handlers"
	"github.com/brave-intl/bat-go/utils/logging"
	"github.com/brave/go-sync/cache"
	syncContext "github.com/brave/go-sync/context"
	"github.com/brave/go-sync/controller"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

var (
	commit    string
	version   string
	buildTime string
)

func setupLogger(ctx context.Context) (context.Context, *zerolog.Logger) {
	ctx = context.WithValue(ctx, appctx.EnvironmentCTXKey, os.Getenv("ENV"))
	ctx = context.WithValue(ctx, appctx.LogLevelCTXKey, zerolog.WarnLevel)
	return logging.SetupLogger(ctx)
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
		r.Use(batware.RequestLogger(logger))
	}

	r.Use(chiware.Timeout(60 * time.Second))
	r.Use(batware.BearerToken)
	r.Use(middleware.CommonResponseHeaders)
	r.Use(middleware.Auth)
	r.Use(middleware.DisabledChain)

	db, err := datastore.NewDynamo()
	if err != nil {
		sentry.CaptureException(err)
		log.Panic().Err(err).Msg("Must be able to init datastore to start")
	}

	redis := cache.NewRedisClient()
	cache := cache.NewCache(cache.NewRedisClientWithPrometheus(redis, "redis"))

	// Provide datastore & cache via context
	ctx = context.WithValue(ctx, syncContext.ContextKeyDatastore, db)
	ctx = context.WithValue(ctx, syncContext.ContextKeyCache, &cache)

	r.Mount("/v2", controller.SyncRouter(
		cache,
		datastore.NewDatastoreWithPrometheus(db, "dynamo")))
	r.Get("/metrics", batware.Metrics())

	log.Info().
		Str("version", version).
		Str("commit", commit).
		Str("buildTime", buildTime).
		Msg("server starting up")
	r.Get("/health-check", handlers.HealthCheckHandler(version, buildTime, commit))

	// Add profiling flag to enable profiling routes.
	if os.Getenv("PPROF_ENABLED") != "" {
		// pprof attaches routes to default serve mux
		// host:6061/debug/pprof/
		go func() {
			log.Error().Err(http.ListenAndServe(":6061", http.DefaultServeMux))
		}()
	}

	return ctx, r
}

// StartServer starts the translate proxy server on port 8195
func StartServer() {
	serverCtx, logger := setupLogger(context.Background())

	// Setup Sentry.
	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:     sentryDsn,
			Release: fmt.Sprintf("go-sync@%s-%s", commit, buildTime),
		})
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
