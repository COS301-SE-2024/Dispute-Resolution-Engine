package models

import (
	"time"
)

type User struct {
	ID                int64      `json:"id" gorm:"primaryKey;autoIncrement;column:id"`                                   //Filled in by API
	FirstName         string     `json:"first_name" gorm:"type:varchar(50);not null;column:first_name"`                  //check
	Surname           string     `json:"surname" gorm:"type:varchar(50);not null;column:surname"`                        //check
	Birthdate         time.Time  `json:"birthdate" gorm:"type:date;not null;column:birthdate"`                           //check
	Nationality       string     `json:"nationality" gorm:"type:varchar(50);not null;column:nationality"`                //check
	Role              string     `json:"role" gorm:"type:varchar(50);not null;column:role"`                              //Filled in by API
	Email             string     `json:"email" gorm:"type:varchar(100);unique;not null;column:email"`                    //check
	PasswordHash      string     `json:"password,omitempty" gorm:"type:varchar(255);not null;column:password_hash"`      //Updated by API
	PhoneNumber       *string    `json:"phone_number,omitempty" gorm:"type:varchar(20);column:phone_number"`             //need
	AddressID         *int64     `json:"address_id,omitempty" gorm:"column:address_id"`                                  //what the fuck
	CreatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`                     //Filled in by API
	UpdatedAt         *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:updated_at"`                     //Filled in by API
	LastLogin         *time.Time `gorm:"type:timestamp;column:last_login"`                                               //Filled in by API
	Status            string     `json:"status" gorm:"type:varchar(20);default:'active';column:status"`                  //Filled in by API
	Gender            string     `json:"gender" gorm:"type:gender_enum;column:gender"`                                   //check
	PreferredLanguage *string    `json:"preferred_language,omitempty" gorm:"type:varchar(50);column:preferred_language"` //worked on
	Timezone          *string    `json:"timezone,omitempty" gorm:"type:varchar(50);column:timezone"`                     //need to be handled by me?
	Salt              string     `gorm:"type:varchar(255);column:salt"`
}

type ArchivedDisputeSummary struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Title        string    `json:"title" gorm:"type:varchar(255);column:title"`
	Summary      string    `json:"summary" gorm:"type:text;column:summary"`
	Category     []string  `json:"category" gorm:"type:varchar(255);column:category"`
	DateFiled    time.Time `json:"date_filled" gorm:"type:timestamp;column:date_filled"`
	DateResolved time.Time `json:"date_resolved" gorm:"type:timestamp;column:date_resolved"`
	Resolution   string    `json:"resolution" gorm:"type:text;column:resolution"`
}

type Event struct {
	Timestamp   string `json:"timestamp"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ArchivedDispute struct {
	ArchivedDisputeSummary
	Events []Event `json:"events"`
}

type dispute_decision string

const (
	Resolved   dispute_decision = "Resolved"
	Unresolved dispute_decision = "Unresolved"
	Settled    dispute_decision = "Settled"
	Refused    dispute_decision = "Refused"
	Withdrawn  dispute_decision = "Withdrawn"
	Transfer   dispute_decision = "Transfer"
	Appeal     dispute_decision = "Appeal"
	Other      dispute_decision = "Other"
)

type Dispute struct {
	ID          *int64           `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	CaseDate    time.Time        `json:"case_date" gorm:"type:date;default:CURRENT_DATE;column:case_date"`
	Workflow    *int64           `json:"workflow" gorm:"column:workflow"`
	Status      string           `json:"status" gorm:"type:dispute_status_enum;default:'Awaiting Respondant';column:status"`
	Title       string           `json:"title" gorm:"type:varchar(255);not null;column:title"`
	Description string           `json:"description" gorm:"type:text;column:description"`
	Complainant int64            `json:"complainant" gorm:"column:complainant"`
	Respondant  *int64           `json:"respondant" gorm:"column:respondant"`
	Resolved    bool             `json:"resolved" gorm:"default:false;column:resolved"`
	Decision    dispute_decision `json:"decision" gorm:"type:dispute_decision_enum;default:'Unresolved';column:decision"`
}

type SortAttribute string

const (
	SortByTitle        SortAttribute = "title"
	SortByDateFiled    SortAttribute = "date_filed"
	SortByDateResolved SortAttribute = "date_resolved"
	SortByTimeTaken    SortAttribute = "time_taken"
)

func (User) TableName() string {
	return "users"
}

type Address struct {
	ID          int64   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Country     *string `json:"country,omitempty" gorm:"type:varchar(64);column:code"`
	CountryName *string `json:"coutry_name,omitempty" gorm:"type:varchar(255);column:country"`
	Province    *string `json:"province,omitempty" gorm:"type:varchar(255);column:province"`
	City        *string `json:"city,omitempty" gorm:"type:varchar(255);column:city"`
	Street3     *string `json:"street3,omitempty" gorm:"type:varchar(255);column:street3"`
	Street2     *string `json:"street2,omitempty" gorm:"type:varchar(255);column:street2"`
	Street      *string `json:"street,omitempty" gorm:"type:varchar(255);column:street"`
	AddressType *string `json:"address_type,omitempty" gorm:"type:address_type_enum;column:address_type"`
}

func (Address) TableName() string {
	return "addresses"
}

type Country struct {
	CountryCode string `json:"country_code,omitempty" gorm:"primaryKey;type:varchar(3);not null;column:country_code"`
	CountryName string `json:"country_name,omitempty" gorm:"type:varchar(255);not null;column:country_name"`
}

func (Country) TableName() string {
	return "countries"
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type DisputeExpert struct {
	Dispute int64 `gorm:"primaryKey;column:dispute;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:DisputeID;references:id"`
	User    int64 `gorm:"primaryKey;column:user;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:id"`
}

func (DisputeExpert) TableName() string {
	return "dispute_experts"
}

type File struct {
	ID       *uint     `json:"id" gorm:"primaryKey;column:id;type:serial;autoIncrement:true"`
	FileName string    `json:"label" gorm:"column:file_name;type:varchar(255);not null"`
	Uploaded time.Time `json:"date_submitted" gorm:"column:uploaded;type:timestamp;default:CURRENT_TIMESTAMP"`
	FilePath string    `json:"url" gorm:"column:file_path;type:varchar(255);not null"`
}

func (File) TableName() string {
	return "files"
}

type DisputeEvidence struct {
	Dispute int64 `gorm:"primaryKey;column:dispute;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:DisputeID;references:id"`
	FileID  int64 `gorm:"primaryKey;column:file_id;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FileID;references:id"`
}

func (DisputeEvidence) TableName() string {
	return "dispute_evidence"
}
