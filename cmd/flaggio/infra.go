package main

import (
	"context"
	"net/url"

	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getRedisClient(uri string) (*redis.Client, error) {
	// parse provided uri
	redisURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	// check if password was provided
	redisPass, hasPass := redisURL.User.Password()

	// redis connection options
	redisOpts := &redis.Options{
		Addr: redisURL.Host,
	}
	if hasPass {
		redisOpts.Password = redisPass
	}

	// create redis client & test connection
	redisClient := redis.NewClient(redisOpts)
	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}
	return redisClient, nil
}

func getMongoDatabase(ctx context.Context, uri string) (*mongo.Database, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return mongoClient.Database("flaggio"), nil // TODO: make configurable
}
