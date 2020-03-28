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
	"github.com/victorkt/flaggio/internal/repository/mongodb"
	"github.com/victorkt/flaggio/internal/server/api"
	"github.com/victorkt/flaggio/internal/server/api/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func startAPI(ctx context.Context, c *cli.Context, logger *logrus.Entry) error {
	logger.Info("starting api server ...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.String("database-uri")))
	if err != nil {
		return err
	}

	db := client.Database("flaggio") // TODO: make configurable
	flgRepo, err := mongodb.NewMongoFlagRepository(ctx, db)
	if err != nil {
		return err
	}
	sgmntRepo, err := mongodb.NewMongoSegmentRepository(ctx, db)
	if err != nil {
		return err
	}

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
	)

	srv := &http.Server{
		Addr: c.String("api-addr"),
		Handler: api.NewServer(
			router,
			service.NewFlagService(flgRepo, sgmntRepo),
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
