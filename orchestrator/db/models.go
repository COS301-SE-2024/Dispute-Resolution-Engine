package db

import (
	"encoding/json"
	"time"
)

// Workflow model
type Workflow struct {
    ID                 uint64          `gorm:"primaryKey;autoIncrement"`
	Name               string          `gorm:"type:varchar(100);not null"`
    Definition         json.RawMessage `gorm:"column:definition;type:jsonb"`
    CreatedAt          time.Time       `gorm:"autoCreateTime"`
    AuthorID           *int64          `gorm:"column:author"`
    Author             *User           `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}


// LabelledWorkflows model
type LabelledWorkflow struct {
	WorkflowID uint64   `gorm:"primaryKey;column:workflow_id"`
	TagID      uint64   `gorm:"primaryKey;column:tag_id"`
	Workflow   Workflow `gorm:"foreignKey:WorkflowID"`
	Tag        Tag      `gorm:"foreignKey:TagID"`
}

type User struct {
	ID                int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName         string     `json:"first_name" gorm:"type:varchar(50);not null"`
	Surname           string     `json:"surname" gorm:"type:varchar(50);not null"`
	Birthdate         time.Time  `json:"birthdate" gorm:"type:date;not null"`
	Nationality       string     `json:"nationality" gorm:"type:varchar(50);not null"`
	Role              string     `json:"role" gorm:"type:varchar(50);not null"`
	Email             string     `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash      string     `json:"password,omitempty" gorm:"type:varchar(255);not null"`
	PhoneNumber       *string    `json:"phone_number,omitempty" gorm:"type:varchar(20)"`
	AddressID         *int64     `json:"address_id,omitempty" gorm:"column:address_id"`
	CreatedAt         time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	LastUpdate         *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	LastLogin         *time.Time `gorm:"type:timestamp"`
	Status            string     `json:"status" gorm:"type:varchar(20);default:'active'"`
	Gender            string     `json:"gender" gorm:"type:gender_enum"`
	PreferredLanguage *string    `json:"preferred_language,omitempty" gorm:"type:varchar(50)"`
	Timezone          *string    `json:"timezone,omitempty" gorm:"type:varchar(50)"`
	Salt              string     `gorm:"type:varchar(255)"`
}

// Tag model
type Tag struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	TagName string `gorm:"type:varchar(100);not null"`
}
func (Tag) TableName() string {
	return "tags"
}

func (Workflow) TableName() string {
	return "workflows"
}

func (LabelledWorkflow) TableName() string {
	return "labelled_workflows"
}

func (User) TableName() string {
	return "users"
}
