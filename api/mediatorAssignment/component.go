package mediatorassignment
// AglorithmComponent struct and interface

type AlgorithmComponent interface {
	CalculateScore() float64
}

type BaseComponent struct {
	Function MathFunctions
	DBScore  DBScoreInput
	Operator ComponentOperator
}

func (b *BaseComponent) CalculateScore() float64 {
	score, err := b.DBScore.GetScoreInput()
	if err != nil {
		return 0
	}

	return b.Function.CalculateScore(score)
}

func (b *BaseComponent) ApplyOperator(value1 float64, value2 float64) float64 {
	return b.Operator.ApplyOperator(value1, value2)
}