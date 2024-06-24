package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func Init() *gorm.DB {
	// Load .env file
	if host := os.Getenv("DATABASE_URL"); host == "" {
		err := godotenv.Load("api.env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Retrieve environment variables
	host := os.Getenv("DATABASE_URL")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	// Check if any environment variable is missing
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("One or more required environment variables are missing")
	}

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	for i:= 0; err != nil && i < 10; i++ {
		log.Fatalf("Failed to connect to the database: %v", err)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		time.Sleep(5 * time.Second)
	}

	// Log successful connection
	log.Println("Connected to the database successfully")

	return db
}
