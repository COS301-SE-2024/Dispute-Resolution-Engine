package mediatorassignment


type ComponentOperator interface {
	ApplyOperator(value1 float64, value2 float64) float64
}

type AddOperator struct {
}

func (a *AddOperator) ApplyOperator(value1 float64, value2 float64) float64 {
	return value1 + value2
}

type SubtractOperator struct {
}

func (s *SubtractOperator) ApplyOperator(value1 float64, value2 float64) float64 {
	return value1 - value2
}

type MultiplyOperator struct {
}

func (m *MultiplyOperator) ApplyOperator(value1 float64, value2 float64) float64 {
	return value1 * value2
}

type DivideOperator struct {
}

func (d *DivideOperator) ApplyOperator(value1 float64, value2 float64) float64 {
	return value1 / value2
}