package mediatorassignment

import "math"

// AglorithmComponent struct and interface

type AlgorithmComponent interface {
}

//logorithmic backoff component
//will provide a score which based on a logorithmic scale of the number of active disputes assigned to the mediator
type ComponentAssignedDisputes struct {
}

//exponential ramp up component
// exponential scoring based on the number of disputes rejected by the mediator since last involved
type ComponentRejectionCount struct {
}

//linear growth component
// linear scoring based on the number of disputes resolved by the mediator since last involved
type ComponentTimeSinceLastDispute struct {
}


//------Mathematical functions for the components

type MathFunctions interface {
	//function to calculate the score
	CalculateScore() float64
}

type BaseFunction struct {
	ApplyCapToValue bool
	Cap float64
	InputValue float64
	MoveYAxis float64
	MoveXAxis float64
}

type Expontential struct {
	BaseExponent float64
	BaseFunction
}

func (e *Expontential) CalculateScore() float64 {
	score := e.MoveYAxis + math.Pow(e.BaseExponent, e.InputValue) + e.MoveXAxis

	if e.ApplyCapToValue {
		if score > e.Cap {
			return e.Cap
		}
	}

	return score
}


type Logarithmic struct {
	BaseFunction
	LogBase float64
}

func (l *Logarithmic) CalculateScore() float64 {
	score := l.MoveYAxis + math.Log(l.InputValue) / math.Log(l.LogBase) + l.MoveXAxis

	if l.ApplyCapToValue {
		if score > l.Cap {
			return l.Cap
		}
	}

	return score
}


type Linear struct {
	BaseFunction
	Multiplier float64
}

func (l *Linear) CalculateScore() float64 {
	score := l.MoveYAxis + l.InputValue * l.Multiplier + l.MoveXAxis

	if l.ApplyCapToValue {
		if score > l.Cap {
			return l.Cap
		}
	}

	return score
}

