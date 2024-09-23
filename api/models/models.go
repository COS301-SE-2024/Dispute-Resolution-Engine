package models

import (
	"encoding/json"
	"time"
)

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
	LastUpdate        *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	LastLogin         *time.Time `gorm:"type:timestamp"`
	Status            string     `json:"status" gorm:"type:varchar(20);default:'active'"`
	Gender            string     `json:"gender" gorm:"type:gender_enum"`
	PreferredLanguage *string    `json:"preferred_language,omitempty" gorm:"type:varchar(50)"`
	Timezone          *string    `json:"timezone,omitempty" gorm:"type:varchar(50)"`
	Salt              string     `gorm:"type:varchar(255)"`
}

type ArchivedDisputeSummary struct {
	ID           int64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Title        string   `json:"title" gorm:"type:varchar(255);column:title"`
	Description  string   `json:"description" gorm:"type:text;column:summary"`
	Summary      string   `json:"summary" gorm:"type:text;column:summary"`
	Category     []string `json:"category" gorm:"type:varchar(255);column:category"`
	DateFiled    string   `json:"date_filled" gorm:"type:timestamp;column:date_filled"`
	DateResolved string   `json:"date_resolved" gorm:"type:timestamp;column:date_resolved"`
	Resolution   string   `json:"resolution" gorm:"type:text;column:resolution"`
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

type DisputeStatus string

const (
	StatusAwaitingRespondant DisputeStatus = "Awaiting Respondant"
	StatusActive             DisputeStatus = "Active"
	StatusReview             DisputeStatus = "Review"
	StatusSettled            DisputeStatus = "Settled"
	StatusRefused            DisputeStatus = "Refused"
	StatusWithdrawn          DisputeStatus = "Withdrawn"
	StatusTransfer           DisputeStatus = "Transfer"
	StatusAppeal             DisputeStatus = "Appeal"
	StatusOther              DisputeStatus = "Other"
)

type Dispute struct {
	ID          *int64        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	CaseDate    time.Time     `json:"case_date" gorm:"type:date;default:CURRENT_DATE;column:case_date"`
	Workflow    *int64        `json:"workflow" gorm:"column:workflow"`
	Status      DisputeStatus `json:"status" gorm:"type:dispute_status_enum;default:'Awaiting Respondant';column:status"`
	Title       string        `json:"title" gorm:"type:varchar(255);not null;column:title"`
	Description string        `json:"description" gorm:"type:text;column:description"`
	Complainant int64         `json:"complainant" gorm:"column:complainant"`
	Respondant  *int64        `json:"respondant" gorm:"column:respondant"`
}

type Workflow struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement"`
	Name       string          `gorm:"type:varchar(100);not null"`
	Definition json.RawMessage `gorm:"column:definition;type:jsonb"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`

	AuthorID int64 `gorm:"column:author"`
	Author   *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

func (Workflow) TableName() string {
	return "workflows"
}

type Tag struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	TagName string `gorm:"type:varchar(100);not null"`
}

func (Tag) TableName() string {
	return "tags"
}

type WorkflowTags struct {
	WorkflowID uint64   `gorm:"primaryKey;column:workflow_id"`
	TagID      uint64   `gorm:"primaryKey;column:tag_id"`
	Workflow   Workflow `gorm:"foreignKey:WorkflowID"`
	Tag        Tag      `gorm:"foreignKey:TagID"`
}

func (WorkflowTags) TableName() string {
	return "workflow_tags"
}

type DisputeSummaries struct {
	ID      int64  `json:"id" gorm:"primaryKey;autoincrement;column:dispute"`
	Summary string `json:"summary" gorm:"type:text;column:summary"`
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

type EventTypes string

const (
	Notification EventTypes = "NOTIFICATION"
	Disputes     EventTypes = "DISPUTE"
	Users        EventTypes = "USER"
	Experts      EventTypes = "EXPERT"
	Workflows    EventTypes = "WORKFLOW"
)

// EventLog represents the event_log table
type EventLog struct {
	ID        uint                   `gorm:"primaryKey"`
	CreatedAt time.Time              `gorm:"default:CURRENT_TIMESTAMP"`
	EventType EventTypes             `gorm:"type:event_types"`
	EventData map[string]interface{} `gorm:"type:json"`
}

func (EventLog) TableName() string {
	return "event_log"
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

type ExpertStatus string

const (
	ExpertApproved ExpertStatus = "Approved"
	ExpertReview   ExpertStatus = "Review"
	ExpertRejected ExpertStatus = "Rejected"
)

type DisputeExpert struct {
	Dispute int64        `gorm:"primaryKey;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;references:id"`
	Expert  int64        `gorm:"primaryKey;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;references:id"`
	Status  ExpertStatus `gorm:"<-:false;type:expert_status`
}

func (DisputeExpert) TableName() string {
	return "dispute_experts_view"
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
	UserID  int64 `gorm:"primaryKey;column:user_id;type:bigint;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:id"`
}

func (DisputeEvidence) TableName() string {
	return "dispute_evidence"
}

type ExpObjStatus string

const (
	ObjectionReview    ExpObjStatus = "Review"
	ObjectionSustained ExpObjStatus = "Sustained"
	ObjectionOverruled ExpObjStatus = "Overruled"
)

type ExpertObjection struct {
	ID        int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	DisputeID int64        `gorm:"not null" json:"dispute_id"`
	ExpertID  int64        `gorm:"not null" json:"expert_id"`
	UserID    int64        `gorm:"not null" json:"user_id"`
	Reason    string       `gorm:"type:text" json:"reason"`
	Status    ExpObjStatus `gorm:"type:exp_obj_status;default:'Review'" json:"status"`
}

func (ExpertObjection) TableName() string {
	return "expert_objections"
}

type ActiveWorkflows struct {
	ID               int64           `gorm:"primaryKey;autoIncrement"`
	Workflow         int64           `gorm:"not null"`                               // Foreign Key to Workflow
	CurrentState     string          `gorm:"column:current_state;type:varchar(255)"` // Current State
	DateSubmitted    time.Time       `gorm:"column:date_submitted;type:timestamp"`   // Date the workflow was submitted
	StateDeadline    time.Time       `gorm:"column:state_deadline;type:timestamp"`   // Deadline for the current state
	WorkflowInstance json.RawMessage `gorm:"type:jsonb"`
}

func (ActiveWorkflows) TableName() string {
	return "active_workflows"
}
