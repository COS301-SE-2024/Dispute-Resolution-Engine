package models

import (
	"time"
)

type UpdateUser struct {
	FirstName          string  `json:"first_name"`
	Surname            string  `json:"surname"`
	Phone_number       *string `json:"phone_number"`
	Gender             string  `json:"gender"`
	Nationality        string  `json:"nationality"`
	Timezone           *string `json:"timezone"`
	Preferred_language *string `json:"preferred_language"`
}

type GetUser struct {
	FirstName   string  `json:"first_name"`
	Surname     string  `json:"surname"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`

	Birthdate   string `json:"birthdate"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`

	Timezone          *string `json:"timezone"`
	PreferredLanguage *string `json:"preferred_language"`

	Address Address `json:"address"`
	Theme   string  `json:"theme"`
}

type ColumnValueComparison struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

type DateRange struct {
	Column    string  `json:"column"`
	StartDate *string `json:"startDate, omitempty"`
	EndDate   *string `json:"endDate, omitempty"`
}
type OrderBy struct {
	Column string `json:"column"`
	Order  string `json:"order"`
}
type UserAnalytics struct {
	ColumnvalueComparisons *[]ColumnValueComparison `json:"columnvalueComparisons,omitempty"`
	OrderBy                *[]OrderBy               `json:"orderBy,omitempty"`
	DateRanges             *[]DateRange             `json:"dateRanges,omitempty"`
	GroupBy                *[]string                `json:"groupBy,omitempty"`
	Count                  bool                     `json:"count,omitempty"`
}

type CreateUser struct {
	//These are all the user details that are required to create a user
	FirstName         string  `json:"first_name"`
	Surname           string  `json:"surname"`
	Birthdate         string  `json:"birthdate"`
	Nationality       string  `json:"nationality"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	PhoneNumber       *string `json:"phone_number"`
	Gender            string  `json:"gender"`
	PreferredLanguage *string `json:"preferred_language"`
	Timezone          *string `json:"timezone"`
	UserType          string  `json:"user_type"`
}

type VerifyUser struct {
	Pin string `json:"pin"`
}

type DeleteUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateAddress struct {
	Country     *string `json:"country"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	Street3     *string `json:"street3"`
	Street2     *string `json:"street2"`
	Street      *string `json:"street"`
	AddressType *string `json:"address_type"`
}

type ArchiveSearchRequest struct {
	Search *string          `json:"search,omitempty"`
	Limit  *int             `json:"limit,omitempty"`
	Offset *int             `json:"offset,omitempty"`
	Order  *string          `json:"order,omitempty"`
	Sort   *SortAttribute   `json:"sort,omitempty"`
	Filter *FilterAttribute `json:"filter,omitempty"`
}

type FilterAttribute struct {
	Category []string `json:"category"`
	Time     []int    `json:"time"`
}

type CreateDispute struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Respondent  Respondent `json:"respondent"`
}

type Respondent struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

type SendResetRequest struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	NewPassword string `json:"newPassword"`
}

type ExpertApproveRequest struct {
	ExpertID string `json:"expert_id"`
}

type ExpertRejectRequest struct {
	ExpertID *int64  `json:"expert_id"`
	Reason   *string `json:"reason"`
}

type RejectExpertReview struct {
	Status *ExpObjStatus  `json:"status"`
}

type DisputeStatusChange struct {
	Status string `json:"status"`
}

type RecommendExpert struct {
	DisputeId int `json:"dispute_id"`
}

type RejectExpert struct {
	DisputeId int64 `json:"dispute_id"`
	ExpertId  int64 `json:"expert_id"`
}

type SortOrder string
type SortAttributeAdmin string
type FilterAttributeAdmin string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"

	SortAttributeTitle        SortAttributeAdmin = "title"
	SortAttributeStatus       SortAttributeAdmin = "status"
	SortAttributeWorkflow     SortAttributeAdmin = "workflow"
	SortAttributeDateFiled    SortAttributeAdmin = "date_filed"
	SortAttributeDateResolved SortAttributeAdmin = "date_resolved"

	FilterAttributeStatus   FilterAttributeAdmin = "status"
	FilterAttributeWorkflow FilterAttributeAdmin = "workflow"
)

type Filter struct {
	// The attribute to filter by
	Attr string `json:"attr"`

	// The value to search for
	Value string `json:"value"`
}

type DateFilter struct {
	Filed    *FiledDate    `json:"filed,omitempty"`
	Resolved *ResolvedDate `json:"resolved,omitempty"`
}

type FiledDate struct {
	Before *string `json:"before,omitempty"`
	After  *string `json:"after,omitempty"`
}

type ResolvedDate struct {
	Before *string `json:"before,omitempty"`
	After  *string `json:"after,omitempty"`
}

type Sort struct {
	Attr  string `json:"attr"`
	Order string `json:"order,omitempty"`
}

type AdminDisputesRequest struct {
	// Search term for the title of disputes
	Search *string `json:"search,omitempty"`

	// Pagination parameters
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`

	Sort *Sort `json:"sort,omitempty"`

	// The filters to apply to data
	Filter []Filter `json:"filter,omitempty"`

	DateFilter *DateFilter `json:"dateFilter,omitempty"`
}

type TicketsRequest struct {
	// Search term for the title of disputes
	Search *string `json:"search,omitempty"`

	// Pagination parameters
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`

	Sort *Sort `json:"sort,omitempty"`

	// The filters to apply to data
	Filter []Filter `json:"filter,omitempty"`
}

type PatchTicketStatus struct {
	Status string `json:"status"`
}

type TicketMessageCreate struct {
	Message string `json:"message"`
}

type TicketCreate struct {
	DisputeID int64  `json:"dispute_id"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type ViewExpetRejectionsRequest struct {
	Expert_id  *int64 `json:"expert_id,omitempty"`
	Dispute_id *int64 `json:"dispute_id,omitempty"`
	Limits     *int   `json:"limits,omitempty"`
	Offset     *int   `json:"offset,omitempty"`
}

type CreateWorkflow struct {
	Name       string               `json:"name,omitempty"`
	Definition WorkflowOrchestrator `json:"definition,omitempty"`
	// Category   []int64                `json:"category,omitempty"`
}

type UpdateWorkflow struct {
	Name               *string               `json:"name,omitempty"`
	WorkflowDefinition *WorkflowOrchestrator `json:"definition,omitempty"`
	// Category           *[]int64                `json:"category,omitempty"`
	// Author             *int64                  `json:"author,omitempty"`
}

type NewActiveWorkflow struct {
	DisputeID *int64 `json:"dispute_id,omitempty"`
	Workflow  *int64 `json:"workflow_id,omitempty"`
}

type ResetActiveWorkflow struct {
	DisputeID    *int64     `json:"dispute_id,omitempty"`
	CurrentState *string    `json:"current_state,omitempty"`
	Deadline     *time.Time `json:"deadline,omitempty"`
}

type GetWorkflow struct {
	Search *string `json:"search,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
}

type NotifyEventOrchestrator struct {
	ActiveWorkflowID *int64  `json:"id"`
	CurrentState     *string `json:"current_state"`
}

type AdminGroupingAnalytics struct {
	Group *string `json:"group"`
}
