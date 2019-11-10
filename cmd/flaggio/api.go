package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger}))
	router.Use(middleware.Recoverer)

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   c.StringSlice("cors-allowed-origins"),
		AllowCredentials: true,
		Debug:            c.Bool("cors-debug"),
	}).Handler)

	port := "8080"
	srv := &http.Server{
		Addr: ":" + port,
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
	logger.Infof("api server started. listening on port %s", port)
	return srv.ListenAndServe()
}
