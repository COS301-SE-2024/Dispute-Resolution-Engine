package redisDB

import (
	"api/utilities"
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	maxRetryAttempts = 5
	retryTimeout     = 1
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
	for i := 0; err != nil && i < maxRetryAttempts; i++ {
		_, err = rdb.Ping(ctx).Result()
		time.Sleep(retryTimeout * time.Second)
	}
	if err != nil {
		return nil, err
	}

	// Log successful connection

	RDB = rdb
	return rdb, nil
}
