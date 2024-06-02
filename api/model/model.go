package model

import (
	"crypto/rand"
	"encoding/json"

	"golang.org/x/crypto/argon2"
)

type BaseRequest struct {
	RequestType string          `json:"request_type"`
	Body        json.RawMessage `json:"body"`
}

type Argon2idHash struct {
	time    uint32
	memory  uint32
	threads uint8
	keylen  uint32
	saltlen uint32
}

type HashSalt struct {
	Hash []byte
	Salt []byte
}

func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keylen uint32) *Argon2idHash {
	return &Argon2idHash{time, memory, threads, keylen, saltLen}
}

func randomSalt(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (a *Argon2idHash) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	if(len(salt) == 0) {
		salt, err = randomSalt(a.saltlen)
	}
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keylen)
	return &HashSalt{Hash: hash, Salt: salt}, nil

}

type CreateAccountBody struct {
	FirstName    string `json:"first_name"`
	Surname      string `json:"surname"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

type LoginBody struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type DisputeSummaryBody struct {
	UserID string `json:"id"`
}

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// CREATE TABLE addresses (
//     id BIGINT PRIMARY KEY,
//     code VARCHAR(64),
//     country VARCHAR(255),
//     province VARCHAR(255),
//     city VARCHAR(255),
//     street3 VARCHAR(255),
//     street2 VARCHAR(255),
//     street VARCHAR(255),
//     address_type INTEGER,
//     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

type Address struct {
	Id           int    `json:"id"`
	Code         string `json:"code"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Street3      string `json:"street3"`
	Street2      string `json:"street2"`
	Street       string `json:"street"`
	Address_type int    `json:"address_type"`
	Last_updated string `json:"last_updated"`
}

// CREATE TYPE gender_enum AS ENUM ('male', 'female', 'non-binary', 'prefer not to say', 'other');

// CREATE TABLE users (
//     id SERIAL PRIMARY KEY,
//     first_name VARCHAR(50) NOT NULL,
//     surname VARCHAR(50) NOT NULL,
//     birthdate DATE NOT NULL,
//     nationality VARCHAR(50) NOT NULL,
//     role VARCHAR(50) NOT NULL,
//     email VARCHAR(100) NOT NULL UNIQUE,
//     password_hash VARCHAR(255) NOT NULL,
//     phone_number VARCHAR(20),
//     address_id BIGINT REFERENCES addresses(id),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     last_login TIMESTAMP,
//     status VARCHAR(20) DEFAULT 'active',
//     gender gender_enum,
//     preferred_language VARCHAR(50),
//     timezone VARCHAR(50)
// );

type User struct {
	ID                 int    `json:"id"`
	First_name         string `json:"first_name"`
	Surname            string `json:"surname"`
	Birthdate          string `json:"birthdate"`
	Nationality        string `json:"national"`
	Role               string `json:"role"`
	Email              string `json:"email"`
	Password_hash      string `json:"password_hash"`
	Phone_number       string `json:"phone_number"`
	Address_id         int    `json:"address_id"`
	Created_at         string `json:"created_at"`
	Updated_at         string `json:"updated_at"`
	Last_login         string `json:"last_login"`
	Status             string `json:"status"`
	Gender             string `json:"gender"`
	Preferred_language string `json:"preferred_language"`
	Timezone           string `json:"timezone"`
}

type DisputeSummary struct {
	DisputeID    string `json:"id"`
	DisputeTitle string `json:"title"`
}

type LoginUser struct {
	Email         string `json:"email"`
	Password_hash string `json:"password_hash"`
}

func AuthUser() *LoginUser {
	return &LoginUser{}
}

func NewUser() *User {
	return &User{}
}

func NewAddress() *Address {
	return &Address{}
}
