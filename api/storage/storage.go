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

func (s *PostgresStore) Init() error {
	var err error

	err = s.createAddressesTable()
	if err != nil {
		return err
	}

	err = s.createGenderEnum()
	if err != nil {
		return err
	}

	err = s.createUsersTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) createAddressesTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS addresses (
        id BIGINT PRIMARY KEY,
        code VARCHAR(64),
        country VARCHAR(255),
        province VARCHAR(255),
        city VARCHAR(255),
        street3 VARCHAR(255),
        street2 VARCHAR(255),
        street VARCHAR(255),
        address_type INTEGER,
        last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
    _, err := s.db.Exec(query)
    if err != nil {
        return fmt.Errorf("could not create addresses table: %w", err)
    }
    return nil
}

func (s *PostgresStore) createGenderEnum() error {
    query := `
    DO $$ BEGIN
        CREATE TYPE gender_enum AS ENUM ('male', 'female', 'non-binary', 'prefer not to say', 'other');
    EXCEPTION
        WHEN duplicate_object THEN NULL;
    END $$;
    `
    _, err := s.db.Exec(query)
    if err != nil {
        return fmt.Errorf("could not create gender_enum type: %w", err)
    }
    return nil
}

func (s *PostgresStore) createUsersTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50) NOT NULL,
        surname VARCHAR(50) NOT NULL,
        birthdate DATE NOT NULL,
        nationality VARCHAR(50) NOT NULL,
        role VARCHAR(50) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        phone_number VARCHAR(20),
        address_id BIGINT REFERENCES addresses(id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        last_login TIMESTAMP,
        status VARCHAR(20) DEFAULT 'active',
        gender gender_enum,
        preferred_language VARCHAR(50),
        timezone VARCHAR(50)
    );
    `
    _, err := s.db.Exec(query)
    if err != nil {
        return fmt.Errorf("could not create users table: %w", err)
    }
    return nil
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

