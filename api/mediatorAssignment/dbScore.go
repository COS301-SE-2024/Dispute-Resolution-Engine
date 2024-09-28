package mediatorassignment

import (
	"api/models"
	"time"
)

type ScoreModeler interface {
	GetScoreInput(summary []models.ExpertSummaryView) []ResultWithID
}

type LastAssignmentstruct struct {
}

func (d *LastAssignmentstruct) GetScoreInput(summary []models.ExpertSummaryView) []ResultWithID {

	score := make([]ResultWithID, len(summary))
	for _, expertSummary := range summary {
		//calculate score last assignment
		lastAssignment := time.Since(expertSummary.LastAssignedDate).Hours() / 24
		score = append(score, ResultWithID{ID: expertSummary.ExpertID, Result: lastAssignment})
	}
	return score
}

type AssignedDisputes struct {
}

func (d *AssignedDisputes) GetScoreInput(summary []models.ExpertSummaryView) []ResultWithID {

	score := make([]ResultWithID, len(summary))
	for _, expertSummary := range summary {
		//calculate score assigned disputes
		score = append(score, ResultWithID{ID: expertSummary.ExpertID, Result: float64(expertSummary.ActiveDisputeCount)})
	}
	return score
}

type RejectionCount struct {
}

func (d *RejectionCount) GetScoreInput(summary []models.ExpertSummaryView) []ResultWithID {

	score := make([]ResultWithID, len(summary))
	for _, expertSummary := range summary {
		//calculate score rejection count
		score = append(score, ResultWithID{ID: expertSummary.ExpertID, Result: expertSummary.RejectionPercentage})
	}
	return score
}
