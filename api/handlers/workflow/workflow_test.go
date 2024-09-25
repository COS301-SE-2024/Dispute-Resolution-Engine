package workflow_test

import (
	"api/models"
	"errors"
)

//------------------------------------------------------------------------- Mocks

//DB Model

type mockDB struct {
	throwError bool
	ReturnWorkflowArray []models.Workflow
	ReturnWorkflow *models.Workflow
	ReturnDispute *models.Dispute
	ReturnTagArray []models.Tag
	ReturnTag *models.Tag
}

func (m *mockDB) GetWorkflowsWithLimitOffset(limit, offset *int) ([]models.Workflow, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return []models.Workflow{}, nil
}

func (m *mockDB) GetWorkflowByID(id uint64) (*models.Workflow, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return &models.Workflow{}, nil
}

func (m *mockDB) FindDipsuteByID(id uint64) (*models.Dispute, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return &models.Dispute{}, nil
}

func (m *mockDB) QueryTagsToRelatedWorkflow(workflowID uint64) ([]models.Tag, error) {
	if m.throwError {
		return nil, errors.New("error")
	}
	return []models.Tag{}, nil
}

func (m *mockDB) CreateWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) CreateWorkflowTags(workflowID uint64, tags []models.Tag) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) CreateActiveWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) UpdateWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) UpdateActiveWorkflow(workflow *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) DeleteTagsByWorkflowID(workflowID uint64) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) DeleteWorkflow(wf *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}

func (m *mockDB) DeleteActiveWorkflow(wf *models.Workflow) error {
	if m.throwError {
		return errors.New("error")
	}
	return nil
}





