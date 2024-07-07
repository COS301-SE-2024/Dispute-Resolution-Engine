package models

type UpdateUser struct {
	FirstName          string  `json:"first_name"`
	Surname            string  `json:"surname"`
	Phone_number       *string `json:"phone_number"`
	Gender             string  `json:"gender"`
	Nationality        string  `json:"nationality"`
	Timezone           *string `json:"timezone"`
	Preferred_language *string `json:"preferred_language"`
}

type GetUser struct {
	FirstName   string  `json:"first_name"`
	Surname     string  `json:"surname"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`

	Birthdate   string `json:"birthdate"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`

	Timezone          *string `json:"timezone"`
	PreferredLanguage *string `json:"preferred_language"`

	Address Address `json:"address"`
	Theme   string  `json:"theme"`
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
}

type VerifyUser struct {
	Pin string `json:"pin"`
}

type DeleteUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateAddress struct {
	Country     *string `json:"country"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	Street3     *string `json:"street3"`
	Street2     *string `json:"street2"`
	Street      *string `json:"street"`
	AddressType *string `json:"address_type"`
}

type ArchiveSearchRequest struct {
	Search *string          `json:"search,omitempty"`
	Limit  *int             `json:"limit,omitempty"`
	Offset *int             `json:"offset,omitempty"`
	Order  *string          `json:"order,omitempty"`
	Sort   *SortAttribute   `json:"sort,omitempty"`
	Filter *FilterAttribute `json:"filter,omitempty"`
}

type FilterAttribute struct {
	Category []string `json:"category"`
	Time     []int    `json:"time"`
}

type CreateDispute struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Respondent  Respondent `json:"respondent"`
}

type Respondent struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type DisputeNotify struct {
	Type      string `json:"type"`
	DisputeID int    `json:"dispute_id"`
}

type SendResetRequest struct {
	Email string `json:"email"`
}
