package api

import (
	"context"
	"log"
	"net/http"

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

func Command() cli.Command {
	return cli.Command{
		Name:  "api",
		Usage: "start the api server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "port",
				Usage:  "Port to listen to",
				EnvVar: "API_PORT",
				Value:  "25880",
			},
			cli.StringFlag{
				Name:   "database-uri",
				Usage:  "Database URI",
				EnvVar: "API_DATABASE_URI",
				Value:  "mongodb://localhost:27017/flaggio",
			},
		},
		Action: run,
	}
}

func run(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(c.String("database-uri")),
	)
	if err != nil {
		return err
	}

	db := client.Database("flaggio") // TODO: make configurable
	flgRepo, err := mongodb.NewMongoFlagRepository(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	sgmntRepo, err := mongodb.NewMongoSegmentRepository(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	resolvers := api.NewResolver(
		flgRepo,
		sgmntRepo,
	)

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Check against your desired domains here
			return r.Host == "example.org"
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(api.NewExecutableSchema(
		api.Config{Resolvers: resolvers}),
		handler.WebsocketUpgrader(upgrader),
	))

	logrus.Infof("connect to http://localhost:%s/ for GraphQL playground", c.String("port"))
	return http.ListenAndServe(":"+c.String("port"), router)
}
