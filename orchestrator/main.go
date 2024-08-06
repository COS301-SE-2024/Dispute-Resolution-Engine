package main

import (
	"fmt"
	"time"
	// "orchestrator/db"
	// "orchestrator/env"
	// "orchestrator/utilities"
	"orchestrator/workflow"
	"orchestrator/statemachine"
)

var RequiredEnvVariables = []string{
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
	period, _ := time.ParseDuration("15s")
	fee_timer2 := workflow.CreateTimer("fee_timer", period, workflow.TriggerFeeNotPaid)
	state2 := workflow.CreateState("state2")
	state2.AddTimer(fee_timer2)
	state3 := workflow.CreateState("state3")
	state4 := workflow.CreateState(workflow.StateDisputeFeeDue)

	wf := workflow.CreateWorkflow(1, "workflow1", state1)
	wf.AddState(state2)
	wf.AddState(state3)
	wf.AddState(state4)

	t1to2 := workflow.CreateTransition("t1to2", state1.GetName(), state2.GetName(), workflow.TriggerResponseReceived)
	t2to3 := workflow.CreateTransition("t2to3", state2.GetName(), state3.GetName(), workflow.TriggerResponseReceived)
	t1to4 := workflow.CreateTransition("t1to4", state1.GetName(), state4.GetName(), workflow.TriggerFeeNotPaid)
	t2to4 := workflow.CreateTransition("t2to4", state2.GetName(), state4.GetName(), workflow.TriggerFeeNotPaid)
	wf.AddTransition(t1to2)
	wf.AddTransition(t2to3)
	wf.AddTransition(t1to4)
	wf.AddTransition(t2to4)

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
	sm := statemachine.NewStateMachine()
	sm.Init(wf)
	sm.Start()
}
