package redisDB

import (
	"api/utilities"
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() (*redis.Client, error) {
	host, err := utilities.GetRequiredEnv("REDIS_URL")
	if err != nil {
		return nil, err
	}

	password, err := utilities.GetRequiredEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}

	db, err := utilities.GetRequiredEnv("REDIS_DB")
	if err != nil {
		return nil, err
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
	return rdb, nil
}
