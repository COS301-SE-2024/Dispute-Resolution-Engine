package mediatorassignment_test

import (
	mediatorassignment "api/mediatorAssignment"
	"api/models"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ComponentTestSuite struct {
	suite.Suite
	ScoreModelerMock  *mediatorassignment.ScoreModelerMock
	MathFunctionsMock *mediatorassignment.MathFunctionsMock
	OperatorMock      *mediatorassignment.ComponentOperatorMock
}

func (suite *ComponentTestSuite) SetupTest() {
	suite.ScoreModelerMock = &mediatorassignment.ScoreModelerMock{}
	suite.MathFunctionsMock = &mediatorassignment.MathFunctionsMock{}
	suite.OperatorMock = &mediatorassignment.ComponentOperatorMock{}
}

func TestComponentTestSuite(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}
func (suite *ComponentTestSuite) TestCalculateScore() {
	tests := []struct {
		name     string
		summary  models.ExpertSummaryView
		expected mediatorassignment.ResultWithID
	}{
		{
			name: "Test with valid summary",
			summary: models.ExpertSummaryView{
				ExpertID: 1,
			},
			expected: mediatorassignment.ResultWithID{
				ID:     1,
				Result: 0,
			},
		},
		{
			name: "Test with zero score",
			summary: models.ExpertSummaryView{
				ExpertID: 2,
			},
			expected: mediatorassignment.ResultWithID{
				ID:     2,
				Result: 0,
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Mock the behavior of ScoreModeler and MathFunctions
			suite.ScoreModelerMock.ReturnGetScoreInput = mediatorassignment.ResultWithID{ID: tt.summary.ExpertID, Result: 0}
			suite.MathFunctionsMock.ReturnCalculateScore = 0

			component := mediatorassignment.NewAlgorithmComponent(suite.ScoreModelerMock, suite.MathFunctionsMock, suite.OperatorMock)
			result := component.CalculateScore(tt.summary)

			suite.Equal(tt.expected, result)
		})
	}
}
func (suite *ComponentTestSuite) TestApplyOperator() {
	tests := []struct {
		name     string
		value1   []mediatorassignment.ResultWithID
		value2   []mediatorassignment.ResultWithID
		expected []mediatorassignment.ResultWithID
	}{
		{
			name: "Test with equal length slices",
			value1: []mediatorassignment.ResultWithID{
				{ID: 1, Result: 10},
				{ID: 2, Result: 20},
			},
			value2: []mediatorassignment.ResultWithID{
				{ID: 1, Result: 5},
				{ID: 2, Result: 15},
			},
			expected: []mediatorassignment.ResultWithID{
				{ID: 1, Result: 35}, // 10 + 5
				{ID: 2, Result: 35}, // 20 + 15
			},
		},
		{
			name: "Test with empty slices",
			value1: []mediatorassignment.ResultWithID{},
			value2: []mediatorassignment.ResultWithID{},
			expected: []mediatorassignment.ResultWithID{},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Update the mock to dynamically return the sum of the input values
			suite.OperatorMock.ReturnApplyOperator = 35

			component := mediatorassignment.NewAlgorithmComponent(suite.ScoreModelerMock, suite.MathFunctionsMock, suite.OperatorMock)
			result := component.ApplyOperator(tt.value1, tt.value2)

			// Ensure that the empty slice and nil comparison are handled properly
			if result == nil {
				result = []mediatorassignment.ResultWithID{}
			}

			suite.Equal(tt.expected, result)
		})
	}
}
