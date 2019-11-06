package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

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

func startAdmin(ctx context.Context, c *cli.Context) (*http.Server, error) {
	logrus.Info("starting admin server ...")
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

	router.Get("/playground", handler.Playground("GraphQL playground", "/query"))
	router.Post("/query", handler.GraphQL(
		admin.NewExecutableSchema(admin.Config{Resolvers: resolvers}),
		handler.WebsocketUpgrader(upgrader),
	))
	if !c.Bool("no-admin-ui") {
		workDir, _ := os.Getwd()
		buildPath := workDir + "/web/build"
		fileServer(router, "/static", http.Dir(buildPath+"/static"))
		fileServer(router, "/images", http.Dir(buildPath+"/images"))
		router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, buildPath+"/index.html")
		})
	}

	port := c.String("admin-port")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		logrus.Infof("admin server started. connect to http://localhost:%s/playground for GraphQL playground", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("admin server failed to listen: %s", err)
		}
	}()
	return srv, nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
