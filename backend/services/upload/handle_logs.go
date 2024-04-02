package main

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func publishLog(log string) {
	PROJECT_ID := os.Getenv("PROJECT_ID")
	publisher := redisClient()
	publisher.Publish(ctx, "logs:"+PROJECT_ID, log)
}
