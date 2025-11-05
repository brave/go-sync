package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // pprof magic
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	appctx "github.com/brave-intl/bat-go/libs/context"
	"github.com/brave-intl/bat-go/libs/handlers"
	"github.com/brave-intl/bat-go/libs/logging"
	batware "github.com/brave-intl/bat-go/libs/middleware"
	sentry "github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	chiware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/brave/go-sync/cache"
	syncContext "github.com/brave/go-sync/context"
	"github.com/brave/go-sync/controller"
	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/middleware"
)

var (
	commit            string
	version           string
	buildTime         string
	healthCheckActive = true
)

type baseCtxFunc = func(net.Listener) context.Context

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

	healthCheckHandler := func(w http.ResponseWriter, r *http.Request) {
		if healthCheckActive {
			handlers.HealthCheckHandler(version, buildTime, commit, nil, nil)(w, r)
		} else {
			w.WriteHeader(http.StatusGone)
		}
	}
	r.Get("/health-check", healthCheckHandler)
	return ctx, r
}

func setupBaseCtx(ctx context.Context) baseCtxFunc {
	return func(_ net.Listener) context.Context {
		return ctx
	}
}

// StartServer starts the sync proxy server on port 8295
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
	srv := http.Server{
		Addr:        port,
		Handler:     r,
		BaseContext: setupBaseCtx(serverCtx),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		<-sig
		log.Info().Msg("SIGTERM received, disabling health check")

		healthCheckActive = false // disable health check

		time.Sleep(60 * time.Second)
		srv.Shutdown(serverCtx)
	}()

	// Add profiling flag to enable profiling routes.
	if on, _ := strconv.ParseBool(os.Getenv("PPROF_ENABLED")); on {
		// pprof attaches routes to default serve mux
		// host:6061/debug/pprof/
		go func() {
			if err := http.ListenAndServe(":6061", http.DefaultServeMux); err != nil {
				log.Err(err).Msg("pprof service returned error")
			}
		}()
	}

	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Info().Msg("HTTP server closed")
	} else if err != nil {
		sentry.CaptureException(err)
		log.Panic().Err(err).Msg("HTTP server start failed!")
	}
}
