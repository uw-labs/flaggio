package main

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newHTTPServer(ctx context.Context, addr string, handler http.Handler, logger *logrus.Entry, wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	wg.Add(1)
	go gracefulServerShutdown(ctx, srv, logger, wg)

	return srv
}

func newRedisClient(ctx context.Context, uri string, logger *logrus.Entry, wg *sync.WaitGroup) (*redis.Client, error) {
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

	wg.Add(1)
	go gracefulRedisClose(ctx, redisClient, logger, wg)

	return redisClient, nil
}

func newMongoDatabase(ctx context.Context, uri string, logger *logrus.Entry, wg *sync.WaitGroup) (*mongo.Database, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	wg.Add(1)
	go gracefulMongoDisconnect(ctx, mongoClient, logger, wg)
	return mongoClient.Database("flaggio"), nil // TODO: make configurable
}

func gracefulMongoDisconnect(ctx context.Context, client *mongo.Client, logger *logrus.Entry, wg *sync.WaitGroup) { // nolint:interfacer // want mongo.Client for consistency
	<-ctx.Done()
	logger.Debug("disconnecting from mongo")
	if err := client.Disconnect(ctx); err != nil {
		logger.WithError(err).Error("failed to disconnect from mongo")
	}
	wg.Done()
}

func gracefulRedisClose(ctx context.Context, client *redis.Client, logger *logrus.Entry, wg *sync.WaitGroup) {
	<-ctx.Done()
	logger.Debug("disconnecting from redis")
	if err := client.Close(); err != nil {
		logger.WithError(err).Error("failed to disconnect from redis")
	}
	wg.Done()
}

func gracefulServerShutdown(ctx context.Context, srv *http.Server, logger *logrus.Entry, wg *sync.WaitGroup) {
	<-ctx.Done()
	logger.Debug("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("could not gracefully shutdown the server")
	}
	wg.Done()
}
