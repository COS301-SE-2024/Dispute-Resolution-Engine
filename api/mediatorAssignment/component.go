package mediatorassignment


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
	CalculateScore() int
	ApplyCap(setCap bool)
	SetCap(value int)
}

type Expontential struct {
	ApplyCapToValue bool
	Cap int
	InputValue int
	Expontent int
}

func (e *Expontential) CalculateScore() int {
	score := e.InputValue * e.Expontent

	if e.ApplyCapToValue {
		if score > e.Cap {
			return e.Cap
		}
	}

	return score
}

func (e *Expontential) ApplyCap(setCap bool) {
	e.ApplyCapToValue = setCap
}

func (e *Expontential) SetCap(value int) {
	e.Cap = value
}

type Logarithmic struct {
	ApplyCap bool
	Cap int
	InputValue int
}



type Linear struct {
	ApplyCap bool
	Cap int
	InputValue int
}
