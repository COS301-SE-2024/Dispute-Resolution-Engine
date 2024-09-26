package models

import "time"

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Total int64       `json:"total,omitempty"`
}

type DisputeSummaryResponse struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      DisputeStatus `json:"status"`
	Role        *string       `json:"role,omitempty"`
}

type Expert struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type DisputeDetailsResponse struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      DisputeStatus `json:"status"`
	DateCreated time.Time     `json:"case_date"`
	Evidence    []Evidence    `json:"evidence"`
	Experts     []Expert      `json:"experts"`
	Role        string        `json:"role"`
}

type Evidence struct {
	ID           uint      `json:"id"`
	FileName     string    `json:"label"`
	Uploaded     time.Time `json:"date_submitted"`
	FilePath     string    `json:"url"`
	UploaderRole string    `json:"uploader_role"`
}

type ArchiveSearchResponse struct {
	Archives []ArchivedDisputeSummary `json:"archives"`
	Total    int64                    `json:"total"`
}

type DisputeCreationResponse struct {
	DisputeID int64 `json:"id"`
}

type WorkflowResp struct {
	Id    int64  `gorm:"column:id"`
	Title string `gorm:"column:name"`
}

type AdminDisputeSummariesResponse struct {
	Id           int64        `json:"id"`
	Title        string       `json:"title"`
	Status       string       `json:"status"`
	Workflow     WorkflowResp `json:"workflow"`
	DateFiled    string       `json:"date_filed"`
	DateResolved *string      `json:"date_resolved,omitempty" gorm:"column:date_resolved"`
}

type GetWorkflowResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created"`
	LastUpdated time.Time `json:"last_updated"`
	Author      AuthorSum `json:"author"`
}

type AuthorSum struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}
type WorkflowResult struct {
	Total int `json:"total"`
	Workflows []GetWorkflowResponse `json:"workflows"`
}

type DetailedWorkflowResponse struct {
	GetWorkflowResponse
	Definition WorkflowOrchestrator `json:"definition"`
}

type WorkflowOrchestrator struct {
	// The ID of the initial state of the workflow
	Initial string `json:"initial"`

	// All the states in the workflow, keyd by their ID
	States map[string]State `json:"states"`
}

type State struct {
	// Human-readable label of the state
	Label string `json:"label"`

	// Human-readable description of the state, describing what the state means and all
	// the triggers that go from this state
	Description string `json:"description"`

	// All the outgoing triggers of the state, keyed by their IDs
	Triggers map[string]Trigger `json:"triggers,omitempty"`

	// The optional timer associated with a state
	Timer *Timer `json:"timer,omitempty"`
}

type Trigger struct {
	// Human-readable label of the trigger
	Label string `json:"label"`

	// The ID of the next state to transition to
	Next string `json:"next_state"`
}

type Timer struct {
	// The duration that the timer should run for
	Duration DurationWrapper `json:"duration"`

	// The ID of the trigger to fire when the timer expires
	OnExpire string `json:"on_expire"`
}

type DurationWrapper struct {
	time.Duration
}