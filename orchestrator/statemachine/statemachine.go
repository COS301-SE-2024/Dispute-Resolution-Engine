package statemachine

import (
	"context"
	"fmt"
	"orchestrator/scheduler"
	"orchestrator/utilities"
	"orchestrator/workflow"
	"strconv"
	"time"

	"github.com/qmuntal/stateless"
)

// Statemachine consumes a workflow and orchestrates the transitions between states

type IStateMachine interface {
	// initialise the state machine with a workflow
	Init(wf_id string, workflow workflow.Workflow, scheduler *scheduler.Scheduler)
	// fire a trigger to transition to the next state
	TriggerTransition(trigger string) error
	// get the current state of the state machine
	GetCurrentState() (string, error)
	// get the deadline of the current state
	GetStateDeadline() (time.Time, error)
}

type StateMachine struct {
	Label                 string
	StatelessStateMachine *stateless.StateMachine
	Workflow              workflow.Workflow
	Scheduler             *scheduler.Scheduler
	api                   workflow.APIWorkflow
}

func NewStateMachine() IStateMachine {
	return &StateMachine{}
}

func (s *StateMachine) Init(wf_id string, wf workflow.Workflow, sch *scheduler.Scheduler) {
	logger := utilities.NewLogger()
	logger.Info("Initialising state machine")

	queryEngine := workflow.CreateWorkflowQuery()
	s.api = workflow.CreateAPIWorkflow(queryEngine)

	// initState := wf.GetInitialState() // this whole sequence is a bit weird, but idk how else to do it
	// initStatePtr := &initState        // without changing the workflow interface
	s.StatelessStateMachine = stateless.NewStateMachine(wf.Initial)
	s.Workflow = wf
	s.Scheduler = sch // 1 second interval

	// For every state in the workflow, add it to the state machine
	for state_id, state := range wf.States {
		// For every related transition from the state, configure the state with the transition
		stateConfig := s.StatelessStateMachine.Configure(state_id)

		for trigger_id, trigger := range state.Triggers {
			stateConfig.Permit(trigger_id, trigger.Next)
		}

		// Configure timer states
		if timer := state.Timer; timer != nil {
			timerName := fmt.Sprintf("%s_%s", wf_id, state_id)

			// If the current state is the initial state, start the timer
			if state_id == wf.Initial {
				s.Scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
					logger.Debug("Timer expired for state", state_id)
					s.StatelessStateMachine.Fire(timer.OnExpire)
					// No need to update the database here, as the state is already the initial state
					// We update it anyway for safety, explanations in next if block

					//check if timer is nil
					new_state := wf.States[state_id].Triggers[timer.OnExpire].Next
					if wf.States[new_state].Timer != nil {
						new_state_deadline := wf.States[new_state].Timer.GetDeadline()
						wf_id_int, err := strconv.Atoi(wf_id)
						if err != nil {
							logger.Error("Error converting wf_id to int", err)
						}
						err = s.api.UpdateActiveWorkflow(wf_id_int, nil, &new_state, nil, &new_state_deadline, nil)
						if err != nil {
							logger.Error("Error updating active_workflow entry in database from timer event", err)
						}
						logger.Info("Sanity update of active_workflow entry in database", wf_id, new_state, new_state_deadline)
						// Send http post containing workflow ID and new state to the api:9000/event endpoint
						awf, err := s.api.FetchActiveWorkflow(int(wf_id_int))
						if err != nil {
							logger.Error("Error fetching active workflow")
							return
						}
				
						label, desc, err := workflow.GetStateDetails(awf.WorkflowInstance, new_state)
						if err != nil {
							logger.Error("Error getting state details")
							return
						}
				
						request := utilities.APIReq{
							ID:           int64(wf_id_int),
							CurrentState: label,
							Description:  desc,
						}
						utilities.APIPostRequest(utilities.API_URL, request)
					} else {
						wf_id_int, err := strconv.Atoi(wf_id)
						if err != nil {
							logger.Error("Error converting wf_id to int", err)
						}
						err = s.api.UpdateActiveWorkflow(wf_id_int, nil, &new_state, nil, nil, nil)
						if err != nil {
							logger.Error("Error updating active_workflow entry in database from timer event", err)
						}
						logger.Info("Sanity update of active_workflow entry in database", wf_id, new_state)
						// Send http post containing workflow ID and new state to the api:9000/event endpoint
						awf, err := s.api.FetchActiveWorkflow(int(wf_id_int))
						if err != nil {
							logger.Error("Error fetching active workflow")
							return
						}
				
						label, desc, err := workflow.GetStateDetails(awf.WorkflowInstance, new_state)
						if err != nil {
							logger.Error("Error getting state details")
							return
						}
				
						request := utilities.APIReq{
							ID:           int64(wf_id_int),
							CurrentState: label,
							Description:  desc,
						}
						utilities.APIPostRequest(utilities.API_URL, request)
					}
				})
			} else {
				// When the state is entered, start the timer
				stateConfig.OnEntry(func(_ context.Context, args ...any) error {
					logger.Debug("New state entered")
					s.Scheduler.AddTimer(timerName, time.Now().Add(timer.GetDuration()), func() {
						logger.Debug("Timer expired for state", state_id)
						s.StatelessStateMachine.Fire(timer.OnExpire)
						// Get the new state ID
						new_state := wf.States[state_id].Triggers[timer.OnExpire].Next
						// Get the deadline that will be set for the new state
						if wf.States[new_state].Timer != nil {
							new_state_deadline := wf.States[new_state].Timer.GetDeadline()
							// Convert wf_id to int
							wf_id_int, err := strconv.Atoi(wf_id)
							if err != nil {
								logger.Error("Error converting wf_id to int", err)
							}
							// Update the active_workflow entry in the database
							err = s.api.UpdateActiveWorkflow(wf_id_int, nil, &new_state, nil, &new_state_deadline, nil)
							if err != nil {
								logger.Error("Error updating active_workflow entry in database", err)
							}
							logger.Info("Updated active_workflow entry in database", wf_id, new_state, new_state_deadline)
							// Send http post containing workflow ID and new state to the api:9000/event endpoint
							awf, err := s.api.FetchActiveWorkflow(int(wf_id_int))
							if err != nil {
								logger.Error("Error fetching active workflow")
								return
							}
					
							label, desc, err := workflow.GetStateDetails(awf.WorkflowInstance, new_state)
							if err != nil {
								logger.Error("Error getting state details")
								return
							}
					
							request := utilities.APIReq{
								ID:           int64(wf_id_int),
								CurrentState: label,
								Description:  desc,
							}
							utilities.APIPostRequest(utilities.API_URL, request)
						} else {
							// Convert wf_id to int
							wf_id_int, err := strconv.Atoi(wf_id)
							if err != nil {
								logger.Error("Error converting wf_id to int", err)
							}
							// Update the active_workflow entry in the database
							err = s.api.UpdateActiveWorkflow(wf_id_int, nil, &new_state, nil, nil, nil)
							if err != nil {
								logger.Error("Error updating active_workflow entry in database", err)
							}
							logger.Info("Updated active_workflow entry in database", wf_id, new_state)
							// Send http post containing workflow ID and new state to the api:9000/event endpoint
							awf, err := s.api.FetchActiveWorkflow(int(wf_id_int))
							if err != nil {
								logger.Error("Error fetching active workflow")
								return
							}
					
							label, desc, err := workflow.GetStateDetails(awf.WorkflowInstance, new_state)
							if err != nil {
								logger.Error("Error getting state details")
								return
							}
					
							request := utilities.APIReq{
								ID:           int64(wf_id_int),
								CurrentState: label,
								Description:  desc,
							}
							utilities.APIPostRequest(utilities.API_URL, request)
						}
					})
					return nil
				})
			}

			// When the state is exited, remove the timer.
			// WARNING: this may cause some kind of race condition when the exit
			stateConfig.OnExit(func(_ context.Context, args ...any) error {
				s.Scheduler.RemoveTimer(timerName)
				return nil
			})
		}
	}
	s.StatelessStateMachine.Fire(wf.Initial)
}

func (s *StateMachine) TriggerTransition(trigger string) error {
	err := s.StatelessStateMachine.Fire(trigger)
	if err != nil {
		return err
	}
	return nil
}

func (s *StateMachine) GetCurrentState() (string, error) {
	stateAny, err := s.StatelessStateMachine.State(context.Background())
	if err != nil {
		return "", err
	}
	state, ok := stateAny.(string)
	if !ok {
		return "", fmt.Errorf("state is not a string")
	}
	return state, nil
}

func (s *StateMachine) GetStateDeadline() (time.Time, error) {
	state, err := s.GetCurrentState()
	if err != nil {
		return time.Time{}, err
	}
	timer := s.Workflow.States[state].Timer
	if timer == nil {
		return time.Time{}, nil // states are allowed to not have timers
	}
	return time.Now().Add(timer.GetDuration()), nil
}
