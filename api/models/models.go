package models

import (
	"encoding/json"
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
	CreatedAt         time.Time  `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`   //Filled in by API
	LastLogin         *time.Time `json:"last_login" gorm:"type:timestamp;column:last_login"`                             //Filled in by API
	Status            string     `json:"status" gorm:"type:varchar(20);default:'active';column:status"`                  //Filled in by API
	Gender            string     `json:"gender" gorm:"type:gender_enum;column:gender"`                                   //check
	PreferredLanguage *string    `json:"preferred_language,omitempty" gorm:"type:varchar(50);column:preferred_language"` //worked on
	Timezone          *string    `json:"timezone,omitempty" gorm:"type:varchar(50);column:timezone"`                     //need to be handled by me?
	Salt              string     `gorm:"type:varchar(255);column:salt"`
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
	ID           *int64        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	CaseDate     time.Time     `json:"case_date" gorm:"type:date;default:CURRENT_DATE;column:case_date"`
	Workflow     int64        `json:"workflow" gorm:"column:workflow"`
	Status       DisputeStatus `json:"status" gorm:"type:dispute_status_enum;default:'Awaiting Respondant';column:status"`
	Title        string        `json:"title" gorm:"type:varchar(255);not null;column:title"`
	Description  string        `json:"description" gorm:"type:text;column:description"`
	Complainant  int64         `json:"complainant" gorm:"column:complainant"`
	Respondant   *int64        `json:"respondant" gorm:"column:respondant"`
	DateResolved *time.Time    `json:"date_resolved" gorm:"type:timestamp;column:date_resolved"`
}

type Workflow struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"type:varchar(100);not null" json:"name"`
	Definition  json.RawMessage `gorm:"column:definition;type:jsonb" json:"definition"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	LastUpdated time.Time       `gorm:"autoUpdateTime" json:"last_updated"`
	AuthorID    int64           `gorm:"column:author" json:"author_id"`
	Author      *User           `gorm:"foreignKey:AuthorID" json:"author,omitempty" json:"author"`
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
	ID       int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	ExpertID int64        `gorm:"not null" json:"expert_id"`
	TicketID int64        `gorm:"not null" json:"ticket_id"`
	Status   ExpObjStatus `gorm:"type:exp_obj_status;default:'Review'" json:"status"`
}

func (ExpertObjection) TableName() string {
	return "expert_objections"
}

type DisputeDecisions struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	DisputeID int64     `json:"dispute_id" gorm:"column:dispute_id"`
	ExpertID  int64     `json:"expert_id" gorm:"column:expert_id"`
	WriteUpID int64     `json:"write_up_id" gorm:"column:writeup_file_id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at"`
}

type Ticket struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	CreatedBy      int64     `gorm:"not null;column:created_by" json:"created_by"`
	DisputeID      int64     `gorm:"not null;column:dispute_id" json:"dispute_id"`
	Subject        string    `gorm:"not null;type:varchar(255);column:subject" json:"subject"`
	Status         string    `gorm:"not null;type:ticket_status_enum;column:status" json:"status"`
	InitialMessage string    `gorm:"type:text;column:initial_message" json:"initial_message"`
}

type TicketMessages struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID    int64     `gorm:"not null;column:user_id" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	Content   string    `gorm:"type:text;not null;column:content" json:"content"`
	FirstName string    `gorm:"type:varchar(50);column:first_name" json:"first_name"`
	Surname   string    `gorm:"type:varchar(50);column:surname" json:"surname"`
}

type ActiveWorkflows struct {
	ID               int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Workflow         int64           `gorm:"not null" json:"workflow"`                                               // Foreign Key to Workflow
	CurrentState     string          `gorm:"column:current_state;type:varchar(255)" json:"current_state"`            // Current State
	DateSubmitted    time.Time       `gorm:"column:date_submitted;type:timestamp" json:"date_submitted"`             // Date the workflow was submitted
	StateDeadline    time.Time       `gorm:"column:state_deadline;type:timestamp" json:"current_deadline,omitempty"` // Deadline for the current state
	WorkflowInstance json.RawMessage `gorm:"type:jsonb" json:"definition"`
}

func (ActiveWorkflows) TableName() string {
	return "active_workflows"
}

// ExpertObjectionsView represents the expert_objections_view SQL view.
type ExpertObjectionsView struct {
	ObjectionID     int       `gorm:"column:objection_id" json:"id"`                  // ID of the objection
	TicketID        int       `gorm:"column:ticket_id" json:"ticket_id"`              // ID of the related ticket
	TicketCreatedAt time.Time `gorm:"column:ticket_created_at" json:"date_submitted"` // Created date of the ticket
	DisputeID       int       `gorm:"column:dispute_id"`                              // ID of the associated dispute
	DisputeTitle    string    `gorm:"column:dispute_title"`                           // Title of the dispute
	ExpertID        int       `gorm:"column:expert_id"`                               // ID of the expert being objected to
	ExpertFullName  string    `gorm:"column:expert_full_name" json:"expert_name"`     // Full name of the expert
	UserID          int       `gorm:"column:user_id"`                                 // ID of the user creating the objection
	UserFullName    string    `gorm:"column:user_full_name" json:"user_name"`         // Full name of the user creating the objection
	ObjectionStatus string    `gorm:"column:objection_status" json:"status"`          // Status of the objection (Review, Sustained, Overruled)
}

// TableName overrides the table name for GORM to map it to the view.
func (ExpertObjectionsView) TableName() string {
	return "expert_objections_view"
}

type ExpertSummaryView struct {
	ExpertID            uint      `gorm:"column:expert_id; primaryKey" json:"expert_id"`
	ExpertName          string    `gorm:"column:expert_name" json:"expert_name"`
	RejectionPercentage float64   `gorm:"column:rejection_percentage" json:"rejection_percentage"`
	LastAssignedDate    time.Time `gorm:"column:last_assigned_date" json:"last_assigned_date"`
	ActiveDisputeCount  int       `gorm:"column:active_dispute_count" json:"active_dispute_count"`
}

// TableName specifies the table name for GORM
func (ExpertSummaryView) TableName() string {
	return "expert_summary_view"
}
