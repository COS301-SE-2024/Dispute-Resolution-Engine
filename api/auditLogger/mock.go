package auditLogger

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