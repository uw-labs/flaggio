package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
	redis_repo "github.com/victorkt/flaggio/internal/repository/redis"
	"github.com/victorkt/flaggio/internal/server/admin"
)

func startAdmin(ctx context.Context, c *cli.Context, logger *logrus.Entry) error {
	logger.Debug("starting admin server ...")
	isDev := c.String("app-env") == "dev"

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
	variantRepo := mongo_repo.NewVariantRepository(flagRepo.(*mongo_repo.FlagRepository))
	ruleRepo := mongo_repo.NewRuleRepository(
		flagRepo.(*mongo_repo.FlagRepository), segmentRepo.(*mongo_repo.SegmentRepository))
	if redisClient != nil {
		flagRepo = redis_repo.NewFlagRepository(redisClient, flagRepo)
		segmentRepo = redis_repo.NewSegmentRepository(redisClient, segmentRepo)
		variantRepo = redis_repo.NewVariantRepository(redisClient, variantRepo, flagRepo)
		ruleRepo = redis_repo.NewRuleRepository(redisClient, ruleRepo, flagRepo, segmentRepo)
	}

	// setup graphql resolver
	resolver := &admin.Resolver{
		FlagRepo:    flagRepo,
		VariantRepo: variantRepo,
		RuleRepo:    ruleRepo,
		SegmentRepo: segmentRepo,
	}

	// setup graphql server
	gqlSrv := handler.New(
		admin.NewExecutableSchema(admin.Config{Resolvers: resolver}),
	)
	gqlSrv.AddTransport(transport.POST{})
	gqlSrv.Use(extension.Introspection{})

	// setup router
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger:  logger,
			NoColor: c.String("log-formatter") != logFormatterText,
		}),
		cors.New(cors.Options{
			AllowedOrigins:   c.StringSlice("cors-allowed-origins"),
			AllowedHeaders:   c.StringSlice("cors-allowed-headers"),
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowCredentials: true,
			Debug:            c.Bool("cors-debug"),
		}).Handler,
	)
	router.Method("POST", "/query", gqlSrv)
	if isDev {
		router.Get("/playground", playground.Handler("GraphQL playground", "/query"))
	}

	// setup admin UI routes
	if !c.Bool("no-admin-ui") {
		workDir, _ := os.Getwd()
		buildPath := workDir + "/web/build"
		if c.IsSet("build-path") {
			buildPath = c.String("build-path")
		}

		fileServer(router, "/static", http.Dir(buildPath+"/static"))
		fileServer(router, "/images", http.Dir(buildPath+"/images"))
		router.Get("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, buildPath+"/manifest.json")
		})
		router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, buildPath+"/index.html")
		})
	}

	// setup http server
	srv := &http.Server{
		Addr:         c.String("admin-addr"),
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Fatalf("admin server shutdown failed: %+v", err)
		}
	}()

	logger.WithFields(logrus.Fields{
		"cache_enabled": c.IsSet("redis-uri"),
		"listening":     c.String("admin-addr"),
		"playground":    isDev,
		"admin_ui":      !c.Bool("no-admin-ui"),
	}).Infof("admin server started")
	return srv.ListenAndServe()
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
