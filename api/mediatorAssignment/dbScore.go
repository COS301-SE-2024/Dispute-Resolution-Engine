package mediatorassignment

import (
	"api/models"
	"time"
)

type ScoreModeler interface {
	GetScoreInput(summary models.ExpertSummaryView) ResultWithID
}

type LastAssignmentstruct struct {
}

func (d *LastAssignmentstruct) GetScoreInput(summary models.ExpertSummaryView) *ResultWithID {
	lastAssignment := time.Since(summary.LastAssignedDate).Hours() / 24
	score := &ResultWithID{ID: summary.ExpertID, Result: lastAssignment}
	return score
}

type AssignedDisputes struct {
}

func (d *AssignedDisputes) GetScoreInput(summary models.ExpertSummaryView) ResultWithID {
	return ResultWithID{ID: summary.ExpertID, Result: float64(summary.ActiveDisputeCount)}
}

type RejectionCount struct {
}

func (d *RejectionCount) GetScoreInput(summary models.ExpertSummaryView) ResultWithID {
	return ResultWithID{ID: summary.ExpertID, Result: float64(summary.RejectionPercentage)}
}
