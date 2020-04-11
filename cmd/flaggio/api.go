package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/victorkt/clientip"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
	redis_repo "github.com/victorkt/flaggio/internal/repository/redis"
	"github.com/victorkt/flaggio/internal/server/api"
	"github.com/victorkt/flaggio/internal/service"
	redis_svc "github.com/victorkt/flaggio/internal/service/redis"
)

func startAPI(ctx context.Context, wg *sync.WaitGroup, logger *logrus.Entry) error {
	logger.Debug("starting api server ...")

	// connect to mongo
	db, err := newMongoDatabase(ctx, cfg.databaseURI, logger, wg)
	if err != nil {
		return err
	}

	var redisClient *redis.Client
	if cfg.isCachingEnabled() {
		// connect to redis
		redisClient, err = newRedisClient(ctx, cfg.redisURI, logger, wg)
		if err != nil {
			return err
		}
	}

	// setup repositories
	flagRepo, err := mongo_repo.NewFlagRepository(ctx, db)
	if err != nil {
		return err
	}
	segmentRepo, err := mongo_repo.NewSegmentRepository(ctx, db)
	if err != nil {
		return err
	}
	if redisClient != nil {
		flagRepo = redis_repo.NewFlagRepository(redisClient, flagRepo)
		segmentRepo = redis_repo.NewSegmentRepository(redisClient, segmentRepo)
	}

	// setup services
	flagService := service.NewFlagService(flagRepo, segmentRepo)
	if redisClient != nil {
		flagService = redis_svc.NewFlagService(redisClient, flagService)
	}

	// setup router
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.Heartbeat("/ready"),
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger:  logger,
			NoColor: cfg.logFormatter != logFormatterText,
		}),
		tracingMiddleware("flaggio-api", logger),
		cors.New(cors.Options{
			AllowedOrigins:   cfg.corsAllowedOrigins.Value(),
			AllowedHeaders:   cfg.corsAllowedHeaders.Value(),
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowCredentials: true,
			Debug:            cfg.corsDebug,
		}).Handler,
		clientip.Middleware,
	)

	// setup API server
	apiSrv := api.NewServer(
		router,
		flagService,
	)

	logger.WithFields(logrus.Fields{
		"caching":   cfg.isCachingEnabled(),
		"tracing":   cfg.isTracingEnabled(),
		"listening": cfg.apiAddr,
	}).Info("api server started")

	// setup http server
	srv := newHTTPServer(ctx, cfg.apiAddr, apiSrv, logger, wg)

	return srv.ListenAndServe()
}
