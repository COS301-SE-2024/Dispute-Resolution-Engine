package models

import (
	"time"
)

type User struct {
	ID                int64      `gorm:"primaryKey;autoIncrement;column:id"`
	FirstName         string     `gorm:"type:varchar(50);not null;column:first_name"`
	Surname           string     `gorm:"type:varchar(50);not null;column:surname"`
	Birthdate         time.Time  `gorm:"type:date;not null;column:birthdate"`
	Nationality       string     `gorm:"type:varchar(50);not null;column:nationality"`
	Role              string     `gorm:"type:varchar(50);not null;column:role"`
	Email             string     `gorm:"type:varchar(100);unique;not null;column:email"`
	PasswordHash      string     `gorm:"type:varchar(255);not null;column:password_hash"`
	PhoneNumber       *string    `gorm:"type:varchar(20);column:phone_number"`
	AddressID         *int64     `gorm:"column:address_id"`
	CreatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:updated_at"`
	LastLogin         *time.Time `gorm:"type:timestamp;column:last_login"`
	Status            string     `gorm:"type:varchar(20);default:'active';column:status"`
	Gender            string     `gorm:"type:gender_enum;column:gender"`
	PreferredLanguage *string    `gorm:"type:varchar(50);column:preferred_language"`
	Timezone          *string    `gorm:"type:varchar(50);column:timezone"`
	Salt              string     `gorm:"type:varchar(255);column:salt"`
}

func (User) TableName() string {
	return "users"
}
