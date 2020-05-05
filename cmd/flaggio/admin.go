package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	mongo_repo "github.com/victorkt/flaggio/internal/repository/mongodb"
	redis_repo "github.com/victorkt/flaggio/internal/repository/redis"
	"github.com/victorkt/flaggio/internal/server/admin"
)

func startAdmin(ctx context.Context, wg *sync.WaitGroup, logger *logrus.Entry) error {
	logger.Debug("starting admin server ...")

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
	evalRepo, err := mongo_repo.NewEvaluationRepository(ctx, db)
	if err != nil {
		return err
	}
	userRepo, err := mongo_repo.NewUserRepository(ctx, db)
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
		evalRepo = redis_repo.NewEvaluationRepository(redisClient, evalRepo)
	}

	// setup graphql resolver
	resolver := &admin.Resolver{
		FlagRepo:       flagRepo,
		VariantRepo:    variantRepo,
		RuleRepo:       ruleRepo,
		SegmentRepo:    segmentRepo,
		EvaluationRepo: evalRepo,
		UserRepo:       userRepo,
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
		middleware.Heartbeat("/ready"),
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger:  logger,
			NoColor: cfg.logFormatter != logFormatterText,
		}),
		tracingMiddleware("flaggio-admin", logger),
		cors.New(cors.Options{
			AllowedOrigins:   cfg.corsAllowedOrigins.Value(),
			AllowedHeaders:   cfg.corsAllowedHeaders.Value(),
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowCredentials: true,
			Debug:            cfg.corsDebug,
		}).Handler,
	)
	router.Method("POST", "/query", gqlSrv)
	if cfg.playgroundEnabled {
		router.Get("/playground", playground.Handler("GraphQL playground", "/query"))
	}

	// setup admin UI routes
	if !cfg.noAdminUI {
		workDir, _ := os.Getwd()
		buildPath := workDir + "/web/build"
		if cfg.uiBuildPath != "" {
			buildPath = cfg.uiBuildPath
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

	logger.WithFields(logrus.Fields{
		"caching":    cfg.isCachingEnabled(),
		"tracing":    cfg.isTracingEnabled(),
		"listening":  cfg.adminAddr,
		"playground": cfg.playgroundEnabled,
		"admin_ui":   !cfg.noAdminUI,
	}).Info("admin server started")

	// setup http server
	srv := newHTTPServer(ctx, cfg.adminAddr, router, logger, wg)

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
