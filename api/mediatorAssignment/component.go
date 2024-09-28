package mediatorassignment

// AglorithmComponent struct and interface

type AlgorithmComponent interface {
}

type BaseComponent struct {
	Function *MathFunctions
	DBScore  *DBScoreInput
	Operator *ComponentOperator
}

// logorithmic backoff component
// will provide a score which based on a logorithmic scale of the number of active disputes assigned to the mediator
type ComponentAssignedDisputes struct {
}

// exponential ramp up component
// exponential scoring based on the number of disputes rejected by the mediator since last involved
type ComponentRejectionCount struct {
}

// linear growth component
// linear scoring based on the number of disputes resolved by the mediator since last involved
type ComponentTimeSinceLastDispute struct {
}
