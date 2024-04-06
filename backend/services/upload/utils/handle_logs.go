package utils

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func PublishLog(log string) {
	var ctx = context.Background()
	PROJECT_ID := os.Getenv("PROJECT_ID")
	publisher := redisClient()
	publisher.Publish(ctx, "logs:"+PROJECT_ID, log)
}
