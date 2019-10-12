package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/victorkohl/flaggio/internal/repository/mongodb"
	"github.com/victorkohl/flaggio/internal/server/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPort = "8081"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:6548/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	db := c.Database("flaggio")
	flgRepo := mongodb.NewMongoFlagRepository(db)
	resolvers := admin.NewResolver(
		flgRepo,
		mongodb.NewMongoVariantRepository(flgRepo),
	)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(admin.NewExecutableSchema(admin.Config{Resolvers: resolvers})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
