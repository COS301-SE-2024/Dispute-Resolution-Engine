package mediatorassignment_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"api/mediatorAssignment"
)

func TestExponentialCalculateScore(t *testing.T) {
	tests := []struct {
		name           string
		exponential    mediatorassignment.Expontential
		inputValue     float64
		expectedResult float64
	}{
		{
			name: "No cap applied",
			exponential: mediatorassignment.Expontential{
				BaseExponent: 2,
				BaseFunction: mediatorassignment.BaseFunction{
					ApplyCapToValue: false,
					MoveYAxis:       1,
					MoveXAxis:       1,
				},
			},
			inputValue:     3,
			expectedResult: 10,
		},
		{
			name: "Cap applied",
			exponential: mediatorassignment.Expontential{
				BaseExponent: 2,
				BaseFunction: mediatorassignment.BaseFunction{
					ApplyCapToValue: true,
					Cap:             5,
					MoveYAxis:       1,
					MoveXAxis:       1,
				},
			},
			inputValue:     3,
			expectedResult: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.exponential.CalculateScore(tt.inputValue)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestLogarithmicCalculateScore(t *testing.T) {
	tests := []struct {
		name           string
		logarithmic    mediatorassignment.Logarithmic
		inputValue     float64
		expectedResult float64
	}{
		{
			name: "No cap applied",
			logarithmic: mediatorassignment.Logarithmic{
				LogBase: 2,
				BaseFunction: mediatorassignment.BaseFunction{
					ApplyCapToValue: false,
					MoveYAxis:       1,
					MoveXAxis:       1,
				},
			},
			inputValue:     8,
			expectedResult: 4.169925, // Adjusted for precision
		},
		{
			name: "Cap applied",
			logarithmic: mediatorassignment.Logarithmic{
				LogBase: 2,
				BaseFunction: mediatorassignment.BaseFunction{
					ApplyCapToValue: true,
					Cap:             3,
					MoveYAxis:       1,
					MoveXAxis:       1,
				},
			},
			inputValue:     8,
			expectedResult: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.logarithmic.CalculateScore(tt.inputValue)
			// Use a small delta for floating-point comparison (tolerance of 0.00001)
			assert.InDelta(t, tt.expectedResult, result, 0.00001, "Expected: %v, Actual: %v", tt.expectedResult, result)
		})
	}
}


func TestLinearCalculateScore(t *testing.T) {
	tests := []struct {
		name           string
		linear         mediatorassignment.Linear
		inputValue     float64
		expectedResult float64
	}{
		{
			name: "No cap applied",
			linear: mediatorassignment.Linear{
				Multiplier: 2,
				BaseFunction: mediatorassignment.BaseFunction{
					ApplyCapToValue: false,
					MoveYAxis:       1,
					MoveXAxis:       1,
				},
			},
			inputValue:     3,
			expectedResult: 8,
		},
		{
			name: "Cap applied",
			linear: mediatorassignment.Linear{
				Multiplier: 2,
				BaseFunction: mediatorassignment.BaseFunction{
					ApplyCapToValue: true,
					Cap:             5,
					MoveYAxis:       1,
					MoveXAxis:       1,
				},
			},
			inputValue:     3,
			expectedResult: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.linear.CalculateScore(tt.inputValue)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}