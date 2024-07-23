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
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DateCreated time.Time `json:"case_date"`
	Evidence    []File    `json:"evidence"`
	Experts     []Expert  `json:"experts"`
	Role        string   `json:"role"`
}

type ArchiveSearchResponse struct {
	Archives []ArchivedDisputeSummary `json:"archives"`
	Total    int64                    `json:"total"`
}
