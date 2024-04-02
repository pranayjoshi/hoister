package main
import "github.com/go-redis/redis/v8"

redisClient := redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

void publishLog(log) {
    redisClient.publish(`logs:${PROJECT_ID}`, JSON.stringify({ log }))
}
