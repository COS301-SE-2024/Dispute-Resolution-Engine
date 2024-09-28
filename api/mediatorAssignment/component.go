package mediatorassignment

import "api/models"

// AglorithmComponent struct and interface

type AlgorithmComponent interface {
	CalculateScore(summary []models.ExpertSummaryView) []ResultWithID
	ApplyOperator(value1 []ResultWithID, value2 []ResultWithID) []ResultWithID
}

type BaseComponent struct {
	ScoreModeler ScoreModeler
	Function MathFunctions
	Operator ComponentOperator
}

func (b *BaseComponent) CalculateScore(summary []models.ExpertSummaryView) []ResultWithID {
	score := b.ScoreModeler.GetScoreInput(summary)
	for _, score := range score {
		score.Result = b.Function.CalculateScore(score.Result)
	}
	return score
}

func (b *BaseComponent) ApplyOperator(value1 []ResultWithID, value2 []ResultWithID) []ResultWithID {
	for i := range value1 {
		value1[i].Result = b.Operator.ApplyOperator(value1[i].Result, value2[i].Result)
	}
	return value1
}


func NewAlgorithmComponent(scoreModeler ScoreModeler, function MathFunctions, operator ComponentOperator) AlgorithmComponent {
	return &BaseComponent{
		ScoreModeler: scoreModeler,
		Function: function,
		Operator: operator,
	}
}
