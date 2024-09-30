package mediatorassignment

import "api/models"

//--------------------------------------- MathFunctionsMock ---------------------------------------
type MathFunctionsMock struct {
	ReturnCalculateScore float64
}

func (m *MathFunctionsMock) CalculateScore(inputValue float64) float64 {
	return m.ReturnCalculateScore
}


//--------------------------------------- ScoreModelerMock ---------------------------------------

type ScoreModelerMock struct {
	ReturnGetScoreInput ResultWithID
}

func (m *ScoreModelerMock) GetScoreInput(summary models.ExpertSummaryView) ResultWithID {
	return m.ReturnGetScoreInput
}

//--------------------------------------- ComponentOperatorMock ---------------------------------------

type ComponentOperatorMock struct {
	ReturnApplyOperator float64
}

func (m *ComponentOperatorMock) ApplyOperator(value1 float64, value2 float64) float64 {
	return m.ReturnApplyOperator
}

//--------------------------------------- AlgorithmComponentMock ---------------------------------------

type AlgorithmComponentMock struct {
	ReturnCalculateScore ResultWithID
	ReturnApplyOperator []ResultWithID
}

func (m *AlgorithmComponentMock) CalculateScore(summary models.ExpertSummaryView) ResultWithID {
	return m.ReturnCalculateScore
}

func (m *AlgorithmComponentMock) ApplyOperator(value1 []ResultWithID, value2 []ResultWithID) []ResultWithID {
	return m.ReturnApplyOperator
}
