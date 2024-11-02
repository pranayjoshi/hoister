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
		REDIS_URL := os.Getenv("REDIS_URL")
		options, err := redis.ParseURL(REDIS_URL)
		if err != nil {
			log.Fatalf("Failed to parse Redis URL: %v", err)
		}
		redisClient = redis.NewClient(options)
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
