package models

import (
	"time"
)

type User struct {
	FirstName         string     `json:"first_name" gorm:"type:varchar(50);not null;column:first_name"`                  //check
	Surname           string     `json:"surname" gorm:"type:varchar(50);not null;column:surname"`                        //check
	Birthdate         time.Time  `json:"birthdate" gorm:"type:date;not null;column:birthdate"`                           //check
	Nationality       string     `json:"nationality" gorm:"type:varchar(50);not null;column:nationality"`                //check
	Role              string     `json:"role" gorm:"type:varchar(50);not null;column:role"`                              //Filled in by API
	Email             string     `json:"email" gorm:"type:varchar(100);unique;not null;column:email"`                    //check
	PasswordHash      string     `json:"password" gorm:"type:varchar(255);not null;column:password_hash"`                //Updated by API
	PhoneNumber       *string    `json:"phone_number,omitempty" gorm:"type:varchar(20);column:phone_number"`             //need
	AddressID         *int64     `json:"address_id,omitempty" gorm:"column:address_id"`                                  //what the fuck
	CreatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`                     //Filled in by API
	UpdatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:updated_at"`                     //Filled in by API
	LastLogin         *time.Time `gorm:"type:timestamp;column:last_login"`                                               //Filled in by API
	Status            string     `json:"status" gorm:"type:varchar(20);default:'active';column:status"`                  //Filled in by API
	Gender            string     `json:"gender" gorm:"type:gender_enum;column:gender"`                                   //check
	PreferredLanguage *string    `json:"preferred_language,omitempty" gorm:"type:varchar(50);column:preferred_language"` //worked on
	Timezone          *string    `json:"timezone,omitempty" gorm:"type:varchar(50);column:timezone"`                     //need to be handled by me?
	Salt              string     `gorm:"type:varchar(255);column:salt"`                                                  //Filled in by API
}

func (User) TableName() string {
	return "users"
}

type Address struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Code        *string    `json:"code,omitempty" gorm:"type:varchar(64);column:code"`
	Country     *string    `json:"country,omitempty" gorm:"type:varchar(255);column:country"`
	Province    *string    `json:"province,omitempty" gorm:"type:varchar(255);column:province"`
	City        *string    `json:"city,omitempty" gorm:"type:varchar(255);column:city"`
	Street3     *string    `json:"street3,omitempty" gorm:"type:varchar(255);column:street3"`
	Street2     *string    `json:"street2,omitempty" gorm:"type:varchar(255);column:street2"`
	Street      *string    `json:"street,omitempty" gorm:"type:varchar(255);column:street"`
	AddressType *int       `json:"address_type,omitempty" gorm:"type:int;column:address_type"`
	LastUpdated *time.Time `json:"last_updated,omitempty" gorm:"type:timestamp without time zone;default:CURRENT_TIMESTAMP;column:last_updated"`
}

func (Address) TableName() string {
	return "addresses"
}

type Country struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	CountryCode string `json:"country_code,omitempty" gorm:"type:varchar(3);not null;column:country_code"`
	CountryName string `json:"country_name,omitempty" gorm:"type:varchar(255);not null;column:country_name"`
}

func (Country) TableName() string {
	return "countries"
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
