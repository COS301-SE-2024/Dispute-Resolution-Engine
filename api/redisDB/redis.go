package redisDB

import (
	"api/env"
	"api/utilities"
	"context"
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
	logger := utilities.NewLogger().LogWithCaller()
	host, err := env.Get("REDIS_URL")
	if err != nil {
		logger.WithError(err).Error("Failed to get REDIS_URL")
		return nil, err
	}

	password, err := env.Get("REDIS_PASSWORD")
	if err != nil {
		logger.WithError(err).Error("Failed to get REDIS_PASSWORD")
		return nil, err
	}

	db, err := env.Get("REDIS_DB")
	if err != nil {
		logger.WithError(err).Error("Failed to get REDIS_DB")
		return nil, err
	}

	// Convert db to integer
	dbNum, err := strconv.Atoi(db)
	if err != nil {
		logger.WithError(err).Fatal("Failed to convert REDIS_DB to integer")
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
