package storage

import (
	"api/model"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"fmt"
	"os"
	//"github.com/joho/godotenv"
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

	//Login
	AuthenticateUser(user *model.LoginUser) error
	GetSalt(user *model.LoginUser) (string, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	var err error
	var db *sql.DB

	// Load .env file
	// err = godotenv.Load(".env")
	// if err != nil {
	// 	return nil, err
	// }

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
	for i := 0; i < 3; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			if i == 2 { // If this was the last attempt, return the error
				return nil, err
			}
			time.Sleep(5 * time.Second) // Wait for 5 seconds before trying again
			continue
		}
	
		// Check if the connection is successful
		err = db.Ping()
		if err != nil {
			if i == 2 { // If this was the last attempt, return the error
				return nil, err
			}
			time.Sleep(5 * time.Second) // Wait for 5 seconds before trying again
			continue
		}
	
		break // If we reached this point, the connection was successful, so we break out of the loop
	}

	fmt.Println("Connected to PostgreSQL database!")

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Close() {
	s.db.Close()
}

func (s *PostgresStore) Ping() error {
	return s.db.Ping()
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

func (s *PostgresStore) GetAllUsers() ([]*model.User, error) {
	query := `SELECT * FROM users;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get all users: %w", err)
	}

	users := make([]*model.User, 0)
	for rows.Next() {
		user := new(model.User)
		err := rows.Scan(
			&user.ID,
			&user.First_name,
			&user.Surname,
			&user.Birthdate,
			&user.Nationality,
			&user.Role,
			&user.Email,
			&user.Password_hash,
			&user.Phone_number,
			&user.Address_id,
			&user.Created_at,
			&user.Updated_at,
			&user.Last_login,
			&user.Status,
			&user.Gender,
			&user.Preferred_language,
			&user.Timezone,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan user: %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over users: %w", err)
	}

	return users, nil
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
	stmt, err := s.db.Prepare(`INSERT INTO users (
		first_name,
		surname,
		birthdate,
		nationality,
		role,
		email,
		password_hash,
		phone_number,
		address_id,
		last_login,
		gender,
		preferred_language,
		timezone,
		salt
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14
		);
		`)
	if err != nil {
		return fmt.Errorf("could not prepare insert statement: %w", err)
	}

	_, err = stmt.Exec(user.First_name,
		user.Surname,
		user.Birthdate,
		user.Nationality,
		user.Role,
		user.Email,
		user.Password_hash,
		user.Phone_number,
		user.Address_id,
		user.Last_login,
		user.Gender,
		user.Preferred_language,
		user.Timezone,
		user.Salt)

	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}
	return nil
}

func (s *PostgresStore) AuthenticateUser(user *model.LoginUser) error {

	result, err := s.CheckUserExists(user)
	if err != nil {
		return fmt.Errorf("could not authenticate user: %w", err)
	}

	if !result {
		return fmt.Errorf("user does not exist")
	}

	stmt, err := s.db.Prepare(`UPDATE users
	SET last_login = CURRENT_TIMESTAMP
	WHERE email = $1 AND password_hash = $2;`)
	if err != nil {
		return fmt.Errorf("could not prepare update statement: %w", err)
	}

	_, err = stmt.Exec(user.Email, user.Password_hash)
	if err != nil {
		return fmt.Errorf("could not authenticate user: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetSalt(user *model.LoginUser) (string, error) {
	query := `SELECT salt FROM users WHERE email = $1;`
	rows, err := s.db.Query(query, user.Email)
	if err != nil {
		return "", fmt.Errorf("could not get user: %w", err)
	}

	var salt string
	for rows.Next() {
		err := rows.Scan(&salt)
		if err != nil {
			return "", fmt.Errorf("could not scan user: %w", err)
		}
	}
	if err = rows.Err(); err != nil {
		return "", fmt.Errorf("could not iterate over users: %w", err)
	}

	return salt, nil
}

func (s *PostgresStore) CheckUserExists(user *model.LoginUser) (bool, error) {
	query := `SELECT * FROM users WHERE email = $1 AND password_hash = $2;`
	rows, err := s.db.Query(query, user.Email, user.Password_hash)
	if err != nil {
		return false, fmt.Errorf("could not get user: %w", err)
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
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
