package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	logger.Info("starting api server ...")
	// connect to mongo
	db, err := getMongoDatabase(ctx, c.String("database-uri"))
	if err != nil {
		return err
	}

	// connect to redis
	redisClient, err := getRedisClient(c.String("redis-uri"))
	if err != nil {
		return err
	}

	// setup repositories
	flgMongoRepo, err := mongo_repo.NewFlagRepository(ctx, db)
	if err != nil {
		return err
	}
	flgRedisRepo := redis_repo.NewFlagRepository(redisClient, flgMongoRepo)
	sgmntMongoRepo, err := mongo_repo.NewSegmentRepository(ctx, db)
	if err != nil {
		return err
	}
	sgmntRedisRepo := redis_repo.NewSegmentRepository(redisClient, sgmntMongoRepo)

	// setup services
	flagService := service.NewFlagService(flgRedisRepo, sgmntRedisRepo)
	flagRedisService := redis_svc.NewFlagService(redisClient, flagService)

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
			flagRedisService,
		),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logrus.Fatalf("api server shutdown failed: %+v", err)
		}
	}()
	logger.WithField("listening", c.String("api-addr")).Infof("api server started")
	return srv.ListenAndServe()
}
