package admin

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/victorkt/flaggio/internal/repository/mongodb"
	"github.com/victorkt/flaggio/internal/server/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Command() cli.Command {
	return cli.Command{
		Name:  "admin",
		Usage: "start the admin server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "port",
				Usage:  "Port to listen to",
				EnvVar: "ADMIN_PORT",
				Value:  "25881",
			},
			cli.StringFlag{
				Name:   "database-uri",
				Usage:  "Database URI",
				EnvVar: "DATABASE_URI",
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
		return err
	}
	sgmntRepo, err := mongodb.NewMongoSegmentRepository(ctx, db)
	if err != nil {
		return err
	}
	resolvers := admin.NewResolver(
		flgRepo,
		mongodb.NewMongoVariantRepository(flgRepo),
		mongodb.NewMongoRuleRepository(flgRepo, sgmntRepo),
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
	router.Handle("/query", handler.GraphQL(
		admin.NewExecutableSchema(admin.Config{Resolvers: resolvers}),
		handler.WebsocketUpgrader(upgrader),
	))

	logrus.Infof("connect to http://localhost:%s/ for GraphQL playground", c.String("port"))
	return http.ListenAndServe(":"+c.String("port"), router)
}
