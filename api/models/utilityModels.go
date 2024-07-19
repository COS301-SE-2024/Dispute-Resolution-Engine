package models

import "time"

type UserInfoJWT struct {
	ID                int64     `json:"id"`
	FirstName         string    `json:"first_name,omitempty"`
	Surname           string    `json:"surname,omitempty"`
	Birthdate         time.Time `json:"birthdate,omitempty"`
	Nationality       string    `json:"nationality,omitempty"`
	Role              string    `json:"role,omitempty"`
	Email             string    `json:"email"`
	PhoneNumber       *string   `json:"phone_number,omitempty"`
	AddressID         *int64    `json:"address_id,omitempty"`
	Status            string    `json:"status,omitempty"`
	Gender            string    `json:"gender,omitempty"`
	PreferredLanguage *string   `json:"preferred_language,omitempty"`
	Timezone          *string   `json:"timezone,omitempty"`
}

type UserVerify struct {
	User
	Pin string `json:"pin"`
}

func ConvertUserToUserVerify(dbUser User, pin string) *UserVerify {
	return &UserVerify{
		User: dbUser,
		Pin:  pin,
	}
}

func ConvertUserVerifyToUser(dbUser UserVerify) *User {
	return &User{
		ID:                dbUser.ID,
		FirstName:         dbUser.FirstName,
		Surname:           dbUser.Surname,
		Birthdate:         dbUser.Birthdate,
		Nationality:       dbUser.Nationality,
		Role:              dbUser.Role,
		Email:             dbUser.Email,
		PhoneNumber:       dbUser.PhoneNumber,
		AddressID:         dbUser.AddressID,
		Status:            dbUser.Status,
		Gender:            dbUser.Gender,
		PreferredLanguage: dbUser.PreferredLanguage,
		Timezone:          dbUser.Timezone,
	}
}

func ConvertUserToJWTUser(dbUser User) *UserInfoJWT {
	return &UserInfoJWT{
		ID:                dbUser.ID,
		FirstName:         dbUser.FirstName,
		Surname:           dbUser.Surname,
		Birthdate:         dbUser.Birthdate,
		Nationality:       dbUser.Nationality,
		Role:              dbUser.Role,
		Email:             dbUser.Email,
		PhoneNumber:       dbUser.PhoneNumber,
		AddressID:         dbUser.AddressID,
		Status:            dbUser.Status,
		Gender:            dbUser.Gender,
		PreferredLanguage: dbUser.PreferredLanguage,
		Timezone:          dbUser.Timezone,
	}
}

type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}
