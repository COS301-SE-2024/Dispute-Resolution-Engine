package mediatorassignment_test

import (
	"testing"
	"time"

	mediatorassignment "api/mediatorAssignment"
	"api/models"

	"github.com/stretchr/testify/assert"
)

func TestLastAssignmentstruct_GetScoreInput(t *testing.T) {
	tests := []struct {
		name           string
		summary        models.ExpertSummaryView
		expectedResult mediatorassignment.ResultWithID
	}{
		{
			name: "Recent assignment",
			summary: models.ExpertSummaryView{
				ExpertID:         1,
				LastAssignedDate: time.Now().AddDate(0, 0, -1), // 1 day ago
			},
			expectedResult: mediatorassignment.ResultWithID{
				ID:     1,
				Result: 1,
			},
		},
		{
			name: "Old assignment",
			summary: models.ExpertSummaryView{
				ExpertID:         2,
				LastAssignedDate: time.Now().AddDate(0, -1, 0), // 1 month ago
			},
			expectedResult: mediatorassignment.ResultWithID{
				ID:     2,
				Result: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mediatorassignment.LastAssignmentstruct{}
			result := d.GetScoreInput(tt.summary)
			assert.InDelta(t, tt.expectedResult.Result, result.Result, 1, "Expected: %v, Actual: %v", tt.expectedResult.Result, result.Result)
			assert.Equal(t, tt.expectedResult.ID, result.ID)
		})
	}
}

func TestAssignedDisputes_GetScoreInput(t *testing.T) {
	tests := []struct {
		name           string
		summary        models.ExpertSummaryView
		expectedResult mediatorassignment.ResultWithID
	}{
		{
			name: "No active disputes",
			summary: models.ExpertSummaryView{
				ExpertID:           1,
				ActiveDisputeCount: 0,
			},
			expectedResult: mediatorassignment.ResultWithID{
				ID:     1,
				Result: 0,
			},
		},
		{
			name: "Some active disputes",
			summary: models.ExpertSummaryView{
				ExpertID:           2,
				ActiveDisputeCount: 5,
			},
			expectedResult: mediatorassignment.ResultWithID{
				ID:     2,
				Result: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mediatorassignment.AssignedDisputes{}
			result := d.GetScoreInput(tt.summary)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestRejectionCount_GetScoreInput(t *testing.T) {
	tests := []struct {
		name           string
		summary        models.ExpertSummaryView
		expectedResult mediatorassignment.ResultWithID
	}{
		{
			name: "No rejections",
			summary: models.ExpertSummaryView{
				ExpertID:            1,
				RejectionPercentage: 0,
			},
			expectedResult: mediatorassignment.ResultWithID{
				ID:     1,
				Result: 0,
			},
		},
		{
			name: "Some rejections",
			summary: models.ExpertSummaryView{
				ExpertID:            2,
				RejectionPercentage: 25,
			},
			expectedResult: mediatorassignment.ResultWithID{
				ID:     2,
				Result: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &mediatorassignment.RejectionCount{}
			result := d.GetScoreInput(tt.summary)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
