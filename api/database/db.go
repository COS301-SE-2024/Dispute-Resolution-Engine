package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var db *sql.DB

func ConnectDB() {
	var err error

	// Load .env file
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// Database credentials
	host := os.Getenv("DATABASE_URL")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	fmt.Println("hello")
	fmt.Println(host, port, user, password, dbname)

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    // Establish connection to PostgreSQL
    db, err = sql.Open("postgres", connStr)
    if err != nil {
		// fmt.Println("Error: ", err)
        panic(err)
    }
    defer db.Close()

    // Check if the connection is successful
    err = db.Ping()
    if err != nil {
		// fmt.Println("Error: ", err)
        panic(err)
    }

    fmt.Println("Connected to PostgreSQL database!")
}
