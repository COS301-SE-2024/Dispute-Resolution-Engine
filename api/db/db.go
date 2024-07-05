package db

import (
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
	host, err := utilities.GetRequiredEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	port, err := utilities.GetRequiredEnv("DATABASE_PORT")
	if err != nil {
		return nil, err
	}

	user, err := utilities.GetRequiredEnv("DATABASE_USER")
	if err != nil {
		return nil, err
	}

	password, err := utilities.GetRequiredEnv("DATABASE_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbname, err := utilities.GetRequiredEnv("DATABASE_NAME")
	if err != nil {
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
		return nil, err
	}

	return db, nil
}
