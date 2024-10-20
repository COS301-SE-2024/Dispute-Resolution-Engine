package auditLogger

import "api/models"

type MockDisputeProceedingsLogger struct {
	ThrowError bool
	Error      error
}

func (m MockDisputeProceedingsLogger) LogDisputeProceedings(proceedingType string, eventData map[string]interface{}) error {
	if m.ThrowError {
		return m.Error
	}
	return nil
}

type MockDBDisputeProceedingsLogger struct {
	ThrowError bool
	Error      error
}

func (m MockDBDisputeProceedingsLogger) CreateLog(log models.EventLog) error {
	if m.ThrowError {
		return m.Error
	}
	return nil
}