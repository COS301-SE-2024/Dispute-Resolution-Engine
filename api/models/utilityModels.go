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
		PasswordHash:      dbUser.PasswordHash,
		PhoneNumber:       dbUser.PhoneNumber,
		AddressID:         dbUser.AddressID,
		CreatedAt:         dbUser.CreatedAt,
		LastUpdate:        dbUser.LastUpdate,
		LastLogin:         dbUser.LastLogin,
		Status:            dbUser.Status,
		Gender:            dbUser.Gender,
		PreferredLanguage: dbUser.PreferredLanguage,
		Timezone:          dbUser.Timezone,
		Salt:              dbUser.Salt,
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

type DefaultUser struct {
	Email     string
	FirstName string
	Surname   string
}

type AdminIntermediate struct {
	Id           int64      `gorm:"column:id"`
	Title        string     `gorm:"column:title"`
	Status       string     `gorm:"column:status"`
	CaseDate     time.Time  `gorm:"column:case_date"`
	DateResolved *time.Time `gorm:"column:date_resolved"`
}

type TicketIntermediate struct {
	Id             int64     `gorm:"column:id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	Subject        string    `gorm:"column:subject"`
	Status         string    `gorm:"column:status"`
	UserID         int64     `gorm:"column:user_id"`
	FirstName      string    `gorm:"column:first_name"`
	Surname        string    `gorm:"column:surname"`
	InitialMessage *string   `gorm:"column:initial_message"`
}

// type TicketMessage struct {
// 	ID
// }
