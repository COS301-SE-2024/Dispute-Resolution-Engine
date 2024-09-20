package models

import "time"

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
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
