package utils

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func redisClient() *redis.Client {
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	return redis.NewClient(&redis.Options{
		Addr: REDIS_ADDR,
	})
}

func PublishLog(log string) {
	var ctx = context.Background()
	PROJECT_ID := os.Getenv("PROJECT_ID")
	publisher := redisClient()
	publisher.Publish(ctx, "logs:"+PROJECT_ID, log)
}
