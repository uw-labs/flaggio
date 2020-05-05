package mongodb_test

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
)

func TestMain(t *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	if err := c.Ping(ctx, nil); err != nil {
		panic(err)
	}
	mongoClient = c
	mongoDB = mongoClient.Database("flaggio_test")
	code := t.Run()
	if err := mongoClient.Disconnect(ctx); err != nil {
		panic(err)
	}
	cancel()
	os.Exit(code)
}
