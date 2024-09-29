package mediatorassignment

import (
	"api/models"
	"api/utilities"
	"math/rand/v2"
	"sort"

	"gorm.io/gorm"
)

// MediatorAssignment struct and interface
type AlgorithmAssignment interface {
	AssignMediator(count, disputeID int) ([]int, error)
	CalculateScore(summaries []models.ExpertSummaryView, componentID int) []ResultWithID
}

type MediatorAssignment struct {
	Components []AlgorithmComponent
	DB         DBModel
}

func (m *MediatorAssignment) AssignMediator(count, disputeID int) ([]int, error) {
	logger := utilities.NewLogger().LogWithCaller()
	// Get all the experts
	experts, err := m.DB.GetExpertSummaryViews()
	if err != nil {
		return nil, err
	}

	logger.Info("Experts\n", experts)

	// Loop through all the experts
	var intermediateResults []ResultWithID
	if len(experts) > 0 {
		logger.Info("Calculating scores for experts")
		intermediateResults = m.CalculateScore(experts, 0)
		for i := 1; i < len(m.Components); i++ {
			logger.Info("results\n", intermediateResults)
			intermediateResults = m.Components[i].ApplyOperator(intermediateResults, m.CalculateScore(experts, i))
		}
	} else {
		intermediateResults = m.assignRandomValues(experts)
	}
	logger.Info("final\n", intermediateResults)


	// Sort the results
	sort.Slice(intermediateResults, func(i, j int) bool {
		return intermediateResults[i].Result > intermediateResults[j].Result
	})

	logger.Info("Sorted results\n", intermediateResults)
	//get top 10 experts and check they are not rejected
	


	topResults, err := m.GetTopResults(intermediateResults, count, disputeID)
	if err != nil {
		return nil, err
	}

	logger.Info("Top results\n", topResults)

	// Get the expert IDs
	var expertIDs []int
	for _, result := range topResults {
		expertIDs = append(expertIDs, int(result.ID))
	}

	logger.Info("Expert IDs\n", expertIDs)
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

func DefaultAlorithmAssignment(db *gorm.DB) *MediatorAssignment {
	dbmodel := &DBModelReal{DB: db}

	return &MediatorAssignment{
		Components: []AlgorithmComponent{
			&BaseComponent{
				ScoreModeler: &LastAssignmentstruct{},
				Function:    &Linear{BaseFunction: BaseFunction{MoveYAxis: 0, MoveXAxis: 0, ApplyCapToValue: true, Cap: 10,}, Multiplier: 1},
				Operator:    &AddOperator{},
			},
			&BaseComponent{
				ScoreModeler: &AssignedDisputes{},
				Function:    &Logarithmic{BaseFunction: BaseFunction{MoveYAxis: 0, MoveXAxis: 0, ApplyCapToValue: true, Cap: 10,}, LogBase: 10},
				Operator:    &AddOperator{},
			},
			&BaseComponent{
				ScoreModeler: &RejectionCount{},
				Function:    &Expontential{BaseFunction: BaseFunction{MoveYAxis: 0, MoveXAxis: 0, ApplyCapToValue: true, Cap: 10,}, BaseExponent: 10},
				Operator:    &AddOperator{},
			},
		},
		DB: dbmodel,
	}
}

func (m *MediatorAssignment) AddComponent(component AlgorithmComponent) {
	m.Components = append(m.Components, component)
}
