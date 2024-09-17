package db

import (
	"encoding/json"
	"time"
)

// Workflow model
type Workflow struct {
	ID                 uint            `gorm:"primaryKey;autoIncrement"`
	WorkflowDefinition json.RawMessage `gorm:"type:json;not null"`
	CreatedAt          time.Time       `gorm:"autoCreateTime"`
	AuthorID           uint            `gorm:"not null"`
	Author             User            `gorm:"foreignKey:AuthorID;references:ID"`
}

// LabelledWorkflows model
type LabelledWorkflows struct {
	WorkflowID uint     `gorm:"primaryKey;not null"`
	TagID      uint     `gorm:"primaryKey;not null"`
	Workflow   Workflow `gorm:"foreignKey:WorkflowID;references:ID"`
	Tag        Tag      `gorm:"foreignKey:TagID;references:ID"`
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
	UpdatedAt         *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	LastLogin         *time.Time `gorm:"type:timestamp"`
	Status            string     `json:"status" gorm:"type:varchar(20);default:'active'"`
	Gender            string     `json:"gender" gorm:"type:gender_enum"`
	PreferredLanguage *string    `json:"preferred_language,omitempty" gorm:"type:varchar(50)"`
	Timezone          *string    `json:"timezone,omitempty" gorm:"type:varchar(50)"`
	Salt              string     `gorm:"type:varchar(255)"`
}

// Tag model
type Tag struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	TagName string `gorm:"type:varchar(255);not null"`
}

func (Tag) TableName() string {
	return "tags"
}

func (Workflow) TableName() string {
	return "workflow"
}

func (LabelledWorkflows) TableName() string {
	return "labelled_workflows"
}

func (User) TableName() string {
	return "users"
}
