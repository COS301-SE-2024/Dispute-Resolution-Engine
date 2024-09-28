package mediatorassignment

import "api/models"

// MediatorAssignment struct and interface
type AlgorithmAssignment interface {
	
}


type MediatorAssignment struct {
	Components []AlgorithmComponent
	DB 	   DBModel
}

func (m *MediatorAssignment) AssignMediator() {
	// Get all the experts
	experts,err := m.DB.GetExpertSummaryViews()
	if err != nil {
		return
	}

	// Loop through all the experts
	if len(experts) > 0 {
		intermediateResults := m.CalculateScore(experts, 0)
		for i := 1; i < len(m.Components); i++ {
			intermediateResults = m.Components[i].ApplyOperator(intermediateResults, m.CalculateScore(experts, i))
		}
	}

	// Sort the results
	
}

func (m *MediatorAssignment) CalculateScore(summaries []models.ExpertSummaryView, componentID int) []ResultWithID {
	results := make([]ResultWithID, len(summaries))
	component := m.Components[componentID]
	for i, summary := range summaries {
		results[i] = component.CalculateScore(summary)
	}
	return results
}

func (m *MediatorAssignment) AddComponent(component AlgorithmComponent) {
	m.Components = append(m.Components, component)
}
