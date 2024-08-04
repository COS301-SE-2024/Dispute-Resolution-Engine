package main

import (
	"fmt"
	// "orchestrator/db"
	// "orchestrator/env"
	// "orchestrator/utilities"
	"orchestrator/workflow"
)

var requiredEnvVariables = []string{
	// PostGres-related variables
	"DATABASE_URL",
	"DATABASE_PORT",
	"DATABASE_USER",
	"DATABASE_PASSWORD",
	"DATABASE_NAME",
}

func main() {
	// fmt.Println("Hello, World!")
	// logger := utilities.NewLogger().LogWithCaller()
	// env.LoadFromFile(".env", "api.env")

	// for _, key := range requiredEnvVariables {
	// 	env.Register(key)
	// }

	// DB, err := db.Init()
	// if err != nil {
	// 	logger.WithError(err).Fatal("Couldn't initialize database connection")
	// }
	// fmt.Println(DB)

	state1 := workflow.CreateState("state1")
	fee_timer := workflow.CreateTimer("fee_timer", 10, workflow.TriggerFeeNotPaid)
	state1.AddTimer(fee_timer)
	state2 := workflow.CreateState("state2")
	state3 := workflow.CreateState("state3")
	wf := workflow.CreateWorkflow(1, "workflow1", state1)

	wf.AddState(state2)
	wf.AddState(state3)

	t1to2 := workflow.CreateTransition("t1to2", state1.GetName(), state2.GetName(), workflow.TriggerResponseReceived)
	t2to3 := workflow.CreateTransition("t2to3", state2.GetName(), state3.GetName(), workflow.TriggerResponseReceived)
	wf.AddTransition(t1to2)
	wf.AddTransition(t2to3)

	fmt.Println(wf.GetID())
	fmt.Println(wf.GetName())
	fmt.Println(wf.GetInitialState())
	states := wf.GetStates()
	for _, s := range states {
		fmt.Println(s)
	}

	transitions := wf.GetTransitions()
	for _, t := range transitions {
		fmt.Println(t)
	}
}
