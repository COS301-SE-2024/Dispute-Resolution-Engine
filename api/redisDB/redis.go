package redisDB

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() *redis.Client {
	// Retrieve environment variables
	host := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	db := os.Getenv("REDIS_DB")

	// Check if any environment variable is missing
	if host == "" || db == "" {
		log.Fatalf("One or more required environment variables are missing")
	}

	// Convert db to integer
	dbNum, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbNum,
	})

	// Ping the Redis server to check connectivity
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Log successful connection
	log.Println("Connected to Redis successfully")

	RDB = rdb
	return rdb
}
