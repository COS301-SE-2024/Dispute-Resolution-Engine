package mediatorassignment

import (
	"api/models"
	"math/rand/v2"
	"sort"
)

// MediatorAssignment struct and interface
type AlgorithmAssignment interface {
}

type MediatorAssignment struct {
	Components []AlgorithmComponent
	DB         DBModel
}

func (m *MediatorAssignment) AssignMediator() ([]int, error) {
	// Get all the experts
	experts, err := m.DB.GetExpertSummaryViews()
	if err != nil {
		return nil, err
	}

	// Loop through all the experts
	var intermediateResults []ResultWithID
	if len(experts) > 0 {
		intermediateResults = m.CalculateScore(experts, 0)
		for i := 1; i < len(m.Components); i++ {
			intermediateResults = m.Components[i].ApplyOperator(intermediateResults, m.CalculateScore(experts, i))
		}
	} else {
		intermediateResults = m.assignRandomValues(experts)
	}

	// Sort the results
	sort.Slice(intermediateResults, func(i, j int) bool {
		return intermediateResults[i].Result > intermediateResults[j].Result
	})

	//get top 10 experts and check they are not rejected
	topResults, err := m.GetTopResults(intermediateResults, 10, 1)
	if err != nil {
		return nil, err
	}

	// Get the expert IDs
	var expertIDs []int
	for _, result := range topResults {
		expertIDs = append(expertIDs, int(result.ID))
	}

	return expertIDs, nil
}

func (m *MediatorAssignment) CalculateScore(summaries []models.ExpertSummaryView, componentID int) []ResultWithID {
	results := make([]ResultWithID, len(summaries))
	component := m.Components[componentID]
	for i, summary := range summaries {
		results[i] = component.CalculateScore(summary)
	}
	return results
}

func (m *MediatorAssignment) assignRandomValues(summaries []models.ExpertSummaryView) []ResultWithID {
	// assign random values to the experts
	results := make([]ResultWithID, len(summaries))
	for i, summary := range summaries {
		results[i] = ResultWithID{ID: summary.ExpertID, Result: rand.Float64()}
	}
	return results
}

func (m *MediatorAssignment) GetTopResults(results []ResultWithID, count int, disputeID int) ([]ResultWithID,error) {
	rejectedExperts, err := m.DB.GetRejectionFromDispute(disputeID)
	if err != nil {
		return nil, err
	}

	var topResults []ResultWithID
	index := 0
	for len(topResults) < count && index < len(results) {
		if !m.isExpertRejected(rejectedExperts, results[index].ID) {
			topResults = append(topResults, results[index])
		}
		index++
	}

	return topResults, nil
}

func (m *MediatorAssignment) isExpertRejected(rejectedExperts []models.DisputeExpert, expertID uint) bool {
	for _, rejectedExpert := range rejectedExperts {
		if rejectedExpert.Expert == int64(expertID) {
			return true
		}
	}
	return false
}

func (m *MediatorAssignment) AddComponent(component AlgorithmComponent) {
	m.Components = append(m.Components, component)
}
