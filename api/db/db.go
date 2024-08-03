package db

import (
	"api/env"
	"api/utilities"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxRetryAttempts = 10
	retryTimeout     = 2
)

func Init() (*gorm.DB, error) {
	logger := utilities.NewLogger().LogWithCaller()
	envLoader := env.NewEnvLoader()
	host, err := envLoader.Get("DATABASE_URL")
	if err != nil {
		logger.WithError(err).Error("Failed to get DATABASE_URL")
		return nil, err
	}

	port, err := envLoader.Get("DATABASE_PORT")
	if err != nil {
		logger.WithError(err).Error("Failed to get DATABASE_PORT")
		return nil, err
	}

	user, err := envLoader.Get("DATABASE_USER")
	if err != nil {
		logger.WithError(err).Error("Failed to get DATABASE_USER")
		return nil, err
	}

	password, err := envLoader.Get("DATABASE_PASSWORD")
	if err != nil {
		logger.WithError(err).Error("Failed to get DATABASE_PASSWORD")
		return nil, err
	}

	dbname, err := envLoader.Get("DATABASE_NAME")
	if err != nil {
		logger.WithError(err).Error("Failed to get DATABASE_NAME")
		return nil, err
	}

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	for i := 0; err != nil && i < maxRetryAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		time.Sleep(retryTimeout * time.Second)
	}
	if err != nil {
		logger.WithError(err).Error("Failed to connect to database")
		return nil, err
	}
	logger.Info("Connected to database")
	return db, nil
}
