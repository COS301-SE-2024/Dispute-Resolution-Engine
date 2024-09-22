package workflow_test

import (
	"encoding/json"
	"errors"
	"orchestrator/db"
	"orchestrator/workflow"
	"time"

	"github.com/stretchr/testify/suite"
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
		ID:           1,
		FirstName:    "Test",
		Surname:      "User",
		Birthdate:    time.Now(),
		Nationality:  "Kenyan",
		Role:         "Admin",
		Email:        "test@test.com",
		PasswordHash: "password",
		CreatedAt:    time.Now(),
		Status:       "active",
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

func (tdb *TestDbPositive) CreateLabbelledWorkflows(labelledWorkflow *db.LabelledWorkflow) error {
	return nil
}

func (tdb *TestDbPositive) SaveWorkflowInstance(workflow *db.Workflowdb) error {
	return nil
}

func (tdb *TestDbPositive) DeleteLabelledWorkflowByWorkflowId(id uint64) error {
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

func (tdb *TestDbPositive) SaveActiveWorkflowInstance(activeWorkflow *db.ActiveWorkflows) error {
	return nil
}

// ---------------------------- model - dbPositive ----------------------------
// ---------------------------- model - dbNegative ----------------------------
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

func (tdb *TestDbNegative) DeleteLabelledWorkflowByWorkflowId(id uint64) error {
	return errors.New("error")
}

func (tdb *TestDbNegative) CreateLabbelledWorkflows(labelledWorkflow *db.LabelledWorkflow) error {
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

func (tdb *TestDbNegative) SaveActiveWorkflowInstance(activeWorkflow *db.ActiveWorkflows) error {
	return errors.New("error")
}

//---------------------------- model - dbNegative ----------------------------

type WorkflowAPITestSuitePositive struct {
	suite.Suite
	dbQuery *TestDbPositive
}

func (suite *WorkflowAPITestSuitePositive) SetupTest() {
	suite.dbQuery = &TestDbPositive{}
}

func (suite *WorkflowAPITestSuitePositive) TestFetchWorkflowQuery_Positive() {
	// Create the APIWorkflow instance using the mock database
	testingWorkflowAPI := workflow.APIWorkflow{
		WfQuery: suite.dbQuery,
	}

	// Call FetchWorkflowQuery
	workflow, err := testingWorkflowAPI.FetchWorkflow(1)

	// Assert that no error occurred and the workflow is not nil
	suite.NoError(err)
	suite.NotNil(workflow)

	// Assert specific fields in the workflow for correctness
	suite.Equal(1, workflow.ID)
	suite.Equal("Test Workflow", workflow.Name)
	suite.Equal(1, workflow.AuthorID)
}

func (suite *WorkflowAPITestSuitePositive) TestStoreWorkflow_Positive() {
	// Create the APIWorkflow instance using the mock database
	testingWorkflowAPI := workflow.APIWorkflow{
		WfQuery: suite.dbQuery,
	}

	//call function
	wf := workflow.Workflow{Initial: "dispute_created"}
	tags := []int64{1}
	err := testingWorkflowAPI.StoreWorkflow("Test Workflow", wf, tags, 1)
	// Assert that no error occurred
	suite.NoError(err)
}

func (suite *WorkflowAPITestSuitePositive) TestUpdateWorkflow_Positive() {
	// Create the APIWorkflow instance using the mock database
	testingWorkflowAPI := workflow.APIWorkflow{
		WfQuery: suite.dbQuery,
	}

	//call function
	wf := workflow.Workflow{Initial: "dispute_created"}
	name := "bogus"
	err := testingWorkflowAPI.UpdateWorkflow(1, &name, &wf, nil, nil)
	// Assert that no error occurred
	suite.NoError(err)
}

func (suite *WorkflowAPITestSuitePositive) TestFetchActiverWorkflows_Positive() {
	// Create the APIWorkflow instance using the mock database
	testingWorkflowAPI := workflow.APIWorkflow{
		WfQuery: suite.dbQuery,
	}

	// Call FetchActiveWorkflows
	workflows, err := testingWorkflowAPI.FetchActiveWorkflows()

	// Assert that no error occurred and the workflows is not nil
	suite.NoError(err)
	suite.NotNil(workflows)

	// Assert specific fields in the workflow for correctness
	suite.Equal(1, workflows[0].ID)
	suite.Equal(1, workflows[0].WorkflowID)
	suite.Equal("new state", workflows[0].CurrentState)
}

func (suite *WorkflowAPITestSuitePositive) TestFetchActiveWorkflow_Positive() {
	// Create the APIWorkflow instance using the mock database
	testingWorkflowAPI := workflow.APIWorkflow{
		WfQuery: suite.dbQuery,
	}

	// Call FetchActiveWorkflow
	workflow, err := testingWorkflowAPI.FetchActiveWorkflow(1)

	// Assert that no error occurred and the workflow is not nil
	suite.NoError(err)
	suite.NotNil(workflow)

	// Assert specific fields in the workflow for correctness
	suite.Equal(1, workflow.ID)
	suite.Equal(1, workflow.WorkflowID)
	suite.Equal("new state", workflow.CurrentState)
}

func (suite *WorkflowAPITestSuitePositive) TestUpdateActiveWorkflow_Positive() {
	
}
