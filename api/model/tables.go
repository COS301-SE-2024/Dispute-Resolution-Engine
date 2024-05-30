package model

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
	id 		 int
	code 	 string
	country  string
	province string
	city 	 string
	street3  string
	street2  string
	street   string
	address_type int
	last_updated string
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
	id 				int
	first_name 		string
	surname 		string
	birthdate 		string
	nationality 	string
	role 			string
	email 			string
	password_hash 	string
	phone_number 	string
	address_id 		int
	created_at 		string
	updated_at 		string
	last_login 		string
	status 			string
}

func NewUser() *User {
	return &User{}
}

func NewAddress() *Address {
	return &Address{}
}
