package redis_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v7"
)

var (
	redisClient *redis.Client
)

func TestMain(t *testing.M) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: func() string {
			host := os.Getenv("REDIS_HOST")
			port := os.Getenv("REDIS_PORT")
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "6379"
			}
			return fmt.Sprintf("%s:%s", host, port)
		}(),
	})
	if err := redisClient.Ping().Err(); err != nil {
		panic(err)
	}
	code := t.Run()
	if err := redisClient.Close(); err != nil {
		panic(err)
	}
	os.Exit(code)
}
