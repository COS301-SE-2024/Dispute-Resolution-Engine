// package model

// import (
// 	"encoding/json"
// )

// type BaseRequest struct {
// 	RequestType string          `json:"request_type"`
// 	Body        json.RawMessage `json:"body"`
// }

// type CreateAccountBody struct {
// 	FirstName    string `json:"first_name"`
// 	Surname      string `json:"surname"`
// 	Password string `json:"password"`
// 	Email        string `json:"email"`
// }

// type LoginBody struct {
// 	Email        string `json:"email"`
// 	Password string `json:"password"`
// }

// type DisputeSummaryBody struct {
// 	UserID string `json:"id"`
// }

// type Response struct {
// 	Status int         `json:"status"`
// 	Data   interface{} `json:"data,omitempty"`
// 	Error  string      `json:"error,omitempty"`
// }

// // CREATE TABLE addresses (
// //     id BIGINT PRIMARY KEY,
// //     code VARCHAR(64),
// //     country VARCHAR(255),
// //     province VARCHAR(255),
// //     city VARCHAR(255),
// //     street3 VARCHAR(255),
// //     street2 VARCHAR(255),
// //     street VARCHAR(255),
// //     address_type INTEGER,
// //     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// // );

// type Address struct {
// 	Id           int    `json:"id"`
// 	Code         string `json:"code"`
// 	Country      string `json:"country"`
// 	Province     string `json:"province"`
// 	City         string `json:"city"`
// 	Street3      string `json:"street3"`
// 	Street2      string `json:"street2"`
// 	Street       string `json:"street"`
// 	Address_type int    `json:"address_type"`
// 	Last_updated string `json:"last_updated"`
// }

// // CREATE TYPE gender_enum AS ENUM ('male', 'female', 'non-binary', 'prefer not to say', 'other');

// // CREATE TABLE users (
// //     id SERIAL PRIMARY KEY,
// //     first_name VARCHAR(50) NOT NULL,
// //     surname VARCHAR(50) NOT NULL,
// //     birthdate DATE NOT NULL,
// //     nationality VARCHAR(50) NOT NULL,
// //     role VARCHAR(50) NOT NULL,
// //     email VARCHAR(100) NOT NULL UNIQUE,
// //     password_hash VARCHAR(255) NOT NULL,
// //     phone_number VARCHAR(20),
// //     address_id BIGINT REFERENCES addresses(id),
// //     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// //     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// //     last_login TIMESTAMP,
// //     status VARCHAR(20) DEFAULT 'active',
// //     gender gender_enum,
// //     preferred_language VARCHAR(50),
// //     timezone VARCHAR(50)
// // );

// type User struct {
// 	ID                 int    `json:"id"`
// 	First_name         string `json:"first_name"`
// 	Surname            string `json:"surname"`
// 	Birthdate          string `json:"birthdate"`
// 	Nationality        string `json:"national"`
// 	Role               string `json:"role"`
// 	Email              string `json:"email"`
// 	Password_hash      string `json:"password_hash"`
// 	Phone_number       string `json:"phone_number"`
// 	Address_id         int    `json:"address_id"`
// 	Created_at         string `json:"created_at"`
// 	Updated_at         string `json:"updated_at"`
// 	Last_login         string `json:"last_login"`
// 	Status             string `json:"status"`
// 	Gender             string `json:"gender"`
// 	Preferred_language string `json:"preferred_language"`
// 	Timezone           string `json:"timezone"`
// 	Salt               string `json:"salt"`
// }

// type DisputeSummary struct {
// 	DisputeID    string `json:"id"`
// 	DisputeTitle string `json:"title"`
// }

// type LoginUser struct {
// 	Email         string `json:"email"`
// 	Password_hash string `json:"password_hash"`
// }

// func AuthUser() *LoginUser {
// 	return &LoginUser{}
// }

// func NewUser() *User {
// 	return &User{}
// }

// func NewAddress() *Address {
// 	return &Address{}
// }
