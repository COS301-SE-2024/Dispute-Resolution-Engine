package workflow_test

import (
	"encoding/json"
	"errors"
	"orchestrator/db"
	"time"
)

//---------------------------- model - dbPositive ----------------------------

type TestDbPositive struct {
}

const (
	TemplateWorkflow = `{
	"label": "Domain Dispute",
	"initial": "dispute_created",
	"states": {
	  "dispute_created": {
		"label": "Dispute Created",
		"description": "The dispute has been created and is awaiting further action.",
		"triggers": {
			"complaint_not_compliant": {
				"label": "Complaint Not Compliant",
				"next_state": "complaint_rectification"
			}
		},
		"timer": {
		  "duration": "10s",
		  "on_expire": "complaint_not_compliant"
		}
	  },
	  "complaint_rectification": {
		"label": "Complaint Rectification",
		"description": "The complainant has been notified that the complaint is not compliant and has 5 days to rectify the complaint.",
		"triggers": {
			"fee_not_paid": {
				"label": "Fee Not Paid",
				"next_state": "dispute_fee_due"
			}
		},
		"timer": {
		  "duration": "10s",
		  "on_expire": "fee_not_paid"
		}
	  }
	}
  }`
)

func (tdb *TestDbPositive) FetchWorkflowQuery(id int) (*db.Workflowdb, error) {
	return &db.Workflowdb{
		ID:         1,
		Definition: json.RawMessage(TemplateWorkflow),
		AuthorID:   1,
		Name:       "Test Workflow",
		CreatedAt:  time.Now(),
	}, nil
}

func (tdb *TestDbPositive) FetchUserQuery(id int64) (*db.User, error) {
	return &db.User{
		ID:                1,
		FirstName:         "Test",
		Surname:           "User",
		Birthdate:         time.Now(),
		Nationality:       "Kenyan",
		Role:              "Admin",
		Email:             "test@test.com",
		PasswordHash:      "password",
		CreatedAt:         time.Now(),
		Status:            "active",
	}, nil
}

func (tdb *TestDbPositive) FetchTagsByID(id int64) (*db.Tag, error) {
	return &db.Tag{
		ID:      1,
		TagName: "Test Tag",
	}, nil
}

func (tdb *TestDbPositive) CreateWorkflows(workflow *db.Workflowdb) error {
	return nil
}

func (tdb *TestDbPositive) CreateLabbelledWorkdlows(labelledWorkflow *db.LabelledWorkflow) error {
	return nil
}

func (tdb *TestDbPositive) SaveWorkflowInstance(workflow *db.Workflowdb) error {
	return nil
}

func (tdb *TestDbPositive) FetchActiveWorkflows() ([]db.ActiveWorkflows, error) {
	return []db.ActiveWorkflows{
		{
			ID:               1,
			WorkflowID:       1,
			CurrentState:     "new state",
			DateSubmitted:    time.Now(),
			StateDeadline:    time.Now(),
			WorkflowInstance: json.RawMessage(TemplateWorkflow),
		},
	}, nil
}

func (tdb *TestDbPositive) FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error) {
	return &db.ActiveWorkflows{
		ID:               1,
		WorkflowID:       1,
		CurrentState:     "new state",
		DateSubmitted:    time.Now(),
		StateDeadline:    time.Now(),
		WorkflowInstance: json.RawMessage(TemplateWorkflow),
	}, nil
}

func (tdb *TestDbPositive) SaveActiveWorkflow(activeWorkflow *db.ActiveWorkflows) error {
	return nil
}

//---------------------------- model - dbPositive ----------------------------
//---------------------------- model - dbNegative ----------------------------
type TestDbNegative struct {
}

func (tdb *TestDbNegative) FetchWorkflowQuery(id int) (*db.Workflowdb, error) {
	return nil, errors.New("error")
}

func (tdb *TestDbNegative) FetchUserQuery(id int64) (*db.User, error) {
	return nil, errors.New("error")
}

func (tdb *TestDbNegative) FetchTagsByID(id int64) (*db.Tag, error) {
	return nil, errors.New("error")
}

func (tdb *TestDbNegative) CreateWorkflows(workflow *db.Workflowdb) error {
	return errors.New("error")
}

func (tdb *TestDbNegative) CreateLabbelledWorkdlows(labelledWorkflow *db.LabelledWorkflow) error {
	return errors.New("error")
}

func (tdb *TestDbNegative) SaveWorkflowInstance(workflow *db.Workflowdb) error {
	return errors.New("error")
}

func (tdb *TestDbNegative) FetchActiveWorkflows() ([]db.ActiveWorkflows, error) {
	return nil, errors.New("error")
}

func (tdb *TestDbNegative) FetchActiveWorkflow(id int) (*db.ActiveWorkflows, error) {
	return nil, errors.New("error")
}

func (tdb *TestDbNegative) SaveActiveWorkflow(activeWorkflow *db.ActiveWorkflows) error {
	return errors.New("error")
}
//---------------------------- model - dbNegative ----------------------------

