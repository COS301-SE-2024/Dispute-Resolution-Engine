package models

type UpdateUser struct {
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`

	Code        *string `json:"code"` //This is the country code
	Country     *string `json:"country"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	Street3     *string `json:"street3"`
	Street2     *string `json:"street2"`
	Street      *string `json:"street"`
	AddressType *string `json:"address_type"`
}

type CreateUser struct {
	//These are all the user details that are required to create a user
	FirstName         string  `json:"first_name"`
	Surname           string  `json:"surname"`
	Birthdate         string  `json:"birthdate"`
	Nationality       string  `json:"nationality"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	PhoneNumber       *string `json:"phone_number"`
	Gender            string  `json:"gender"`
	PreferredLanguage *string `json:"preferred_language"`
	Timezone          *string `json:"timezone"`

	//These are the user's address details
	Country     *string `json:"country"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	Street3     *string `json:"street3"`
	Street2     *string `json:"street2"`
	Street      *string `json:"street"`
	AddressType *string `json:"address_type"`
}

type VerifyUser struct {
	Pin string `json:"pin"`
}

type DeleteUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}