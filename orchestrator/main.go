package main

import (
	"fmt"
	"time"

	// "orchestrator/db"
	// "orchestrator/env"
	// "orchestrator/utilities"
	// "orchestrator/env"
	// "orchestrator/scheduler"
	// "orchestrator/statemachine"
	"orchestrator/workflow"
	"orchestrator/controller"
)

var RequiredEnvVariables = []string{
	// PostGres-related variables
	// "DATABASE_URL",
	// "DATABASE_PORT",
	// "DATABASE_USER",
	// "DATABASE_PASSWORD",
	// "DATABASE_NAME",
	"ORCHESTRATOR_KEY",
}

func main() {

	// fmt.Println("Hello, World!")
	// logger := utilities.NewLogger().LogWithCaller()
	// env.LoadFromFile(".env", "api.env")

	// for _, key := range RequiredEnvVariables {
	// 	env.Register(key)
	// }

	// DB, err := db.Init()
	// if err != nil {
	// 	logger.WithError(err).Fatal("Couldn't initialize database connection")
	// }
	// fmt.Println(DB)
	period, _ := time.ParseDuration("15s")
	fee_timer1 := workflow.CreateTimer("fee_timer_1", period, workflow.TriggerFeeNotPaid)
	fee_timer2 := workflow.CreateTimer("fee_timer_2", period, workflow.TriggerFeeNotPaid)
	fee_timer3 := workflow.CreateTimer("fee_timer_3", period, workflow.TriggerFeeNotPaid)

	state1 := workflow.CreateState("state1")
	state1.AddTimer(fee_timer1)
	state2 := workflow.CreateState("state2")
	state2.AddTimer(fee_timer2)
	state3 := workflow.CreateState("state3")
	state3.AddTimer(fee_timer3)
	state4 := workflow.CreateState(workflow.StateDisputeFeeDue)

	wf := workflow.CreateWorkflow(1, "workflow1", state1)
	wf.AddState(state2)
	wf.AddState(state3)
	wf.AddState(state4)

	t1to2 := workflow.CreateTransition("t1to2", state1.GetName(), state2.GetName(), workflow.TriggerResponseReceived)
	t2to3 := workflow.CreateTransition("t2to3", state2.GetName(), state3.GetName(), workflow.TriggerResponseReceived)
	t3to4 := workflow.CreateTransition("t3to4", state3.GetName(), state4.GetName(), workflow.TriggerResponseReceived)
	t1to4 := workflow.CreateTransition("t1to4", state1.GetName(), state4.GetName(), workflow.TriggerFeeNotPaid)
	t2to4 := workflow.CreateTransition("t2to4", state2.GetName(), state4.GetName(), workflow.TriggerFeeNotPaid)
	wf.AddTransition(t1to2)
	wf.AddTransition(t2to3)
	wf.AddTransition(t3to4)
	wf.AddTransition(t1to4)
	wf.AddTransition(t2to4)

	//test the WorkFlowToJSON function
	// jsonStr, err := testWorkFlowToJson(wf)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("Workflow JSON:\n", jsonStr+"\n------------\n")

	//test storing workflow in database
	// var category []int64
	// err = workflow.StoreWorkflowToAPI("http://localhost:8080/workflows", wf, category, nil)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// test the JSONToWorkFlow function
	// err = testJsonToWorkFlow(jsonStr)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

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
	ctrl := controller.NewController()
	ctrl.Start()
	ctrl.RegisterStateMachine(wf)
	ctrl.WaitForSignal() 
	
	//test fetch workflow from database
	// wf2, err := workflow.FetchWorkflowFromAPI("http://localhost:8080/workflows/1")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("====================================")
	// fmt.Println(wf2.GetID())
	// fmt.Println(wf2.GetName())
	// fmt.Println(wf2.GetInitialState())
	// states = wf2.GetStates()
	// for _, s := range states {
	// 	fmt.Println(s)
	// }

	// transitions = wf2.GetTransitions()
	// for _, t := range transitions {
	// 	fmt.Println(t)
	// }
	// fmt.Println("====================================")

	// //test update workflow in database
	// err = workflow.UpdateWorkflowToAPI("http://localhost:8080/workflows/6", wf2, nil, nil)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	

}

func testWorkFlowToJson(wf workflow.IWorkflow) (string, error) {
	// Convert the workflow to JSON string
	jsonStr, err := workflow.WorkFlowToJSON(wf.(*workflow.Workflow))
	if err != nil {
		return "", err
	}
	return jsonStr, nil
}

func testJsonToWorkFlow(jsonStr string) error {
	// Convert the JSON string back to a workflow
	fmt.Println("Converting JSON to Workflow====================")
	wf, err := workflow.JSONToWorkFlow(jsonStr)
	if err != nil {
		return err
	}

	// Print the workflow details
	fmt.Println("Workflow ID:", wf.GetID())
	fmt.Println("Workflow Name:", wf.GetName())
	fmt.Println("Initial State:", wf.GetInitialState())

	fmt.Println("States:")
	states := wf.GetStates()
	for _, s := range states {
		fmt.Println(" - State Name:", s.GetName())
		timers := s.GetTimers()
		for _, t := range timers {
			fmt.Println("   - Timer Name:", t.GetName())
			fmt.Println("   - Duration:", t.GetDuration())
			fmt.Println("   - Will Trigger:", t.WillTrigger())
		}
	}

	fmt.Println("Transitions:")
	transitions := wf.GetTransitions()
	for _, t := range transitions {
		fmt.Println(" - Transition Name:", t.GetName())
		fmt.Println("   - From:", t.GetFrom())
		fmt.Println("   - To:", t.GetTo())
		fmt.Println("   - Trigger:", t.GetTrigger())
	}
	fmt.Println("Workflow JSON conversion successful====================")

	return nil
}
