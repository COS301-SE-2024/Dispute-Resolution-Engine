package db

import (
	"api/utilities"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	for i := 0; err != nil && i < 10; i++ {
		log.Fatalf("Failed to connect to the database: %v", err)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		time.Sleep(5 * time.Second)
	}

	// Log successful connection
	log.Println("Connected to the database successfully")

	return db, nil
}
