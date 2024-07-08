package models

import "time"

type DisputeSummaryResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type DisputeDetailsResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DateCreated time.Time `json:"case_date"`
	Evidence    []string  `json:"evidence"`
	Experts     []string  `json:"experts"`
}
