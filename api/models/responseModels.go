package models

import "time"

type DisputeSummaryResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Role        *string `json:"role,omitempty"`
}

type Expert struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type DisputeDetailsResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DateCreated time.Time  `json:"case_date"`
	Evidence    []Evidence `json:"evidence"`
	Experts     []Expert   `json:"experts"`
	Role        string     `json:"role"`
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
