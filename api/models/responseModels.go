package models

type DisputeSummaryResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type DisputeDetailsResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	DateCreated string   `json:"date_created"`
	Evidence    []string `json:"evidence"`
	Experts     []string `json:"experts"`
}
