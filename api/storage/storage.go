package storage

import (
	"api/model"
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"os"
	"fmt"
)

type Storage interface {
	// Create
	CreateUser(user *model.User) error
	CreateAddress(address *model.Address) error

	// Read
	GetUser(id int) (*model.User, error)
	GetAddress(id int) (*model.Address, error)

	// Update
	UpdateUser(user *model.User) error
	UpdateAddress(address *model.Address) error

	// Delete
	DeleteUser(id int) error
	DeleteAddress(id int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	var err error
	var db *sql.DB

	// Load .env file
	err = godotenv.Load(".env")
	if err != nil {
		return nil, err
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
        return nil, err
    }
    defer db.Close()

    // Check if the connection is successful
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    fmt.Println("Connected to PostgreSQL database!")

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateUser(user *model.User) error {
	return nil
}

func (s *PostgresStore) CreateAddress(address *model.Address) error {
	return nil
}

func (s *PostgresStore) GetUser(id int) (*model.User, error) {
	return nil, nil
}

func (s *PostgresStore) GetAddress(id int) (*model.Address, error) {

	return nil, nil
}

func (s *PostgresStore) UpdateUser(user *model.User) error {
	return nil
}

func (s *PostgresStore) UpdateAddress(address *model.Address) error {
	return nil
}

func (s *PostgresStore) DeleteUser(id int) error {
	return nil
}

func (s *PostgresStore) DeleteAddress(id int) error {
	return nil
}

