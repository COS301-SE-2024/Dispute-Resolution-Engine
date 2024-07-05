package models

import "time"

type UserInfoJWT struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`                                   //Filled in by API
	FirstName         string    `json:"first_name" gorm:"type:varchar(50);not null;column:first_name"`                  //check
	Surname           string    `json:"surname" gorm:"type:varchar(50);not null;column:surname"`                        //check
	Birthdate         time.Time `json:"birthdate" gorm:"type:date;not null;column:birthdate"`                           //check
	Nationality       string    `json:"nationality" gorm:"type:varchar(50);not null;column:nationality"`                //check
	Role              string    `json:"role" gorm:"type:varchar(50);not null;column:role"`                              //Filled in by API
	Email             string    `json:"email" gorm:"type:varchar(100);unique;not null;column:email"`                    //check
	PhoneNumber       *string   `json:"phone_number,omitempty" gorm:"type:varchar(20);column:phone_number"`             //need
	AddressID         *int64    `json:"address_id,omitempty" gorm:"column:address_id"`                                  //what the fuck
	Status            string    `json:"status" gorm:"type:varchar(20);default:'active';column:status"`                  //Filled in by API
	Gender            string    `json:"gender" gorm:"type:gender_enum;column:gender"`                                   //check
	PreferredLanguage *string   `json:"preferred_language,omitempty" gorm:"type:varchar(50);column:preferred_language"` //worked on
	Timezone          *string   `json:"timezone,omitempty" gorm:"type:varchar(50);column:timezone"`                     //need to be handled by me?
}

func ConvertUserToJWTUser(dbUser User) *UserInfoJWT{
	return &UserInfoJWT{
		ID: dbUser.ID,
		FirstName: dbUser.FirstName,
		Surname: dbUser.Surname,
		Birthdate: dbUser.Birthdate,
		Nationality: dbUser.Nationality,
		Role: dbUser.Role,
		Email: dbUser.Email,
		PhoneNumber: dbUser.PhoneNumber,
		AddressID: dbUser.AddressID,
		Status: dbUser.Status,
		Gender: dbUser.Gender,
		PreferredLanguage: dbUser.PreferredLanguage,
		Timezone: dbUser.Timezone,
	}
}