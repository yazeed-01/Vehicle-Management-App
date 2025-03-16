package initializers

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func ConnectToRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // use default DB
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
}
