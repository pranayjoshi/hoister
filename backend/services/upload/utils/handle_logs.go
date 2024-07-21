package utils

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func getRedisClient() *redis.Client {
	once.Do(func() {
		REDIS_ADDR := os.Getenv("REDIS_ADDR")
		redisClient = redis.NewClient(&redis.Options{
			Addr: REDIS_ADDR,
		})
	})
	return redisClient
}

func PublishLog(projectID, logs string) {
	publisher := getRedisClient()
	if err := publisher.Publish(context.Background(), "logs:"+projectID, logs).Err(); err != nil {
		log.Printf("Failed to publish log: %v", err)
	} else {
		log.Println("Log published successfully")
	}
}
