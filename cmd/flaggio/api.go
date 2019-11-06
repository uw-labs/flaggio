package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/victorkt/flaggio/internal/repository/mongodb"
	"github.com/victorkt/flaggio/internal/server/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func startAPI(ctx context.Context, c *cli.Context) (*http.Server, error) {
	logrus.Info("starting API server ...")
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(c.String("database-uri")),
	)
	if err != nil {
		return nil, err
	}

	db := client.Database("flaggio") // TODO: make configurable
	flgRepo, err := mongodb.NewMongoFlagRepository(ctx, db)
	if err != nil {
		return nil, err
	}
	sgmntRepo, err := mongodb.NewMongoSegmentRepository(ctx, db)
	if err != nil {
		return nil, err
	}
	resolvers := api.NewResolver(
		flgRepo,
		sgmntRepo,
	)

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // TODO: make configurable
		AllowCredentials: true,
		// Debug:            true,
	}).Handler)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Check against your desired domains here
			return r.Host == "example.org"
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router.Handle("/playground", handler.Playground("GraphQL playground", "/query")) // TODO: make configurable
	router.Handle("/query", handler.GraphQL(api.NewExecutableSchema(
		api.Config{Resolvers: resolvers}),
		handler.WebsocketUpgrader(upgrader),
	))

	port := c.String("api-port")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		logrus.Infof("api server started. connect to http://localhost:%s/playground for GraphQL playground", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("api server failed to listen: %s", err)
		}
	}()
	return srv, nil
}
