package mediatorassignment

import "math"

type MathFunctions interface {
	//function to calculate the score
	CalculateScore(inputValue float64) float64
}

type BaseFunction struct {
	ApplyCapToValue bool
	Cap             float64
	MoveYAxis       float64
	MoveXAxis       float64
}

type Expontential struct {
	BaseExponent float64
	BaseFunction
}

func (e *Expontential) CalculateScore(inputValue float64) float64 {
	score := e.MoveYAxis + math.Pow(e.BaseExponent, inputValue) + e.MoveXAxis

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

func (l *Logarithmic) CalculateScore(inputValue float64) float64 {
	score := l.MoveYAxis + math.Log(inputValue+l.MoveXAxis)/math.Log(l.LogBase)

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

func (l *Linear) CalculateScore(inputValue float64) float64 {
	score := l.MoveYAxis + inputValue*l.Multiplier + l.MoveXAxis

	if l.ApplyCapToValue {
		if score > l.Cap {
			return l.Cap
		}
	}

	return score
}