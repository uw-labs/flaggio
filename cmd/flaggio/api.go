package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/victorkt/clientip"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
	redis_repo "github.com/victorkt/flaggio/internal/repository/redis"
	"github.com/victorkt/flaggio/internal/server/api"
	"github.com/victorkt/flaggio/internal/service"
	redis_svc "github.com/victorkt/flaggio/internal/service/redis"
)

func startAPI(ctx context.Context, c *cli.Context, logger *logrus.Entry) error {
	logger.Debug("starting api server ...")

	// connect to mongo
	db, err := getMongoDatabase(ctx, c.String("database-uri"))
	if err != nil {
		return err
	}

	var redisClient *redis.Client
	if c.IsSet("redis-uri") {
		// connect to redis
		redisClient, err = getRedisClient(c.String("redis-uri"))
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
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger:  logger,
			NoColor: c.String("log-formatter") != "text",
		}),
		cors.New(cors.Options{
			AllowedOrigins:   c.StringSlice("cors-allowed-origins"),
			AllowedHeaders:   c.StringSlice("cors-allowed-headers"),
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowCredentials: true,
			Debug:            c.Bool("cors-debug"),
		}).Handler,
		clientip.Middleware,
	)

	// setup http server
	srv := &http.Server{
		Addr: c.String("api-addr"),
		Handler: api.NewServer(
			router,
			flagService,
		),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Fatalf("api server shutdown failed: %+v", err)
		}
	}()

	logger.WithFields(logrus.Fields{
		"cache_enabled": c.IsSet("redis-uri"),
		"listening":     c.String("api-addr"),
	}).Infof("api server started")
	return srv.ListenAndServe()
}
