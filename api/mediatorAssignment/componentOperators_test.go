package mediatorassignment_test

import (
	"testing"

	mediatorassignment "api/mediatorAssignment"

	"github.com/stretchr/testify/assert"
)

func TestAddOperator(t *testing.T) {
	addOperator := &mediatorassignment.AddOperator{}
	result := addOperator.ApplyOperator(3, 4)
	expected := 7.0
	assert.Equal(t, expected, result)
}

func TestSubtractOperator(t *testing.T) {
	subtractOperator := &mediatorassignment.SubtractOperator{}
	result := subtractOperator.ApplyOperator(10, 4)
	expected := 6.0
	assert.Equal(t, expected, result)
}

func TestMultiplyOperator(t *testing.T) {
	multiplyOperator := &mediatorassignment.MultiplyOperator{}
	result := multiplyOperator.ApplyOperator(3, 4)
	expected := 12.0
	assert.Equal(t, expected, result)
}

func TestDivideOperator(t *testing.T) {
	divideOperator := &mediatorassignment.DivideOperator{}
	result := divideOperator.ApplyOperator(12, 4)
	expected := 3.0
	assert.Equal(t, expected, result)
}

func TestDivideOperator_DivideByZero(t *testing.T) {
	divideOperator := &mediatorassignment.DivideOperator{}
	result := divideOperator.ApplyOperator(12, 0)
	expected := float64(0) // Assuming the function returns 0 for divide by zero
	assert.Equal(t, expected, result)
}
