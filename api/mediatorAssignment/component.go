package mediatorassignment

import "api/models"

// AglorithmComponent struct and interface

type AlgorithmComponent interface {
	CalculateScore(summary models.ExpertSummaryView) ResultWithID
	ApplyOperator(value1 []ResultWithID, value2 []ResultWithID) []ResultWithID 
}

type BaseComponent struct {
	ScoreModeler ScoreModeler
	Function MathFunctions
	Operator ComponentOperator
}

func (b *BaseComponent) CalculateScore(summary models.ExpertSummaryView) ResultWithID {
	result := b.ScoreModeler.GetScoreInput(summary)
	result.Result = b.Function.CalculateScore(result.Result)
	return result
}

func (b *BaseComponent) ApplyOperator(value1 []ResultWithID, value2 []ResultWithID) []ResultWithID {
	var results []ResultWithID
	for i := 0; i < len(value1); i++ {
		results = append(results, ResultWithID{ID: value1[i].ID, Result: b.Operator.ApplyOperator(value1[i].Result, value2[i].Result)})
	}
	return results
}


func NewAlgorithmComponent(scoreModeler ScoreModeler, function MathFunctions, operator ComponentOperator) AlgorithmComponent {
	return &BaseComponent{
		ScoreModeler: scoreModeler,
		Function: function,
		Operator: operator,
	}
}
