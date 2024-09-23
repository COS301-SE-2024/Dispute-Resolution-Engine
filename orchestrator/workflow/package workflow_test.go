package workflow_test

// import (
// 	"errors"
// 	"testing"

// 	"orchestrator/workflow"
// 	"orchestrator/db"

// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/gorm"
// )

// // Mock the db.Init function to return a mock DB instance
// func mockInitSuccess() (*gorm.DB, error) {
// 	return &gorm.DB{}, nil
// }

// // Mock the db.Init function to return an error
// func mockInitError() (*gorm.DB, error) {
// 	return nil, errors.New("db init failed")
// }

// func TestCreateWorkflowQuery(t *testing.T) {
// 	// Mock the db.Init function to return a successful mock DB instance
// 	db.Init = mockInitSuccess

// 	wfq := workflow.CreateWorkflowQuery()

// 	// Check that the WorkflowQuery and its DB instance are not nil
// 	assert.NotNil(t, wfq, "The WorkflowQuery instance should not be nil")
// 	assert.NotNil(t, wfq.DB, "The DB instance in WorkflowQuery should not be nil")
// }

// func TestCreateWorkflowQuery_Error(t *testing.T) {
// 	// Mock the db.Init function to return an error
// 	db.Init = mockInitError

// 	wfq := workflow.CreateWorkflowQuery()

// 	// Check that the WorkflowQuery is nil when db.Init returns an error
// 	assert.Nil(t, wfq, "The WorkflowQuery instance should be nil when db.Init returns an error")
// }
