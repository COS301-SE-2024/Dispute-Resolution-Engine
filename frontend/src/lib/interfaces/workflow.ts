export interface WorkflowSummary {
  id: number;
  name: string;
}

export interface Workflow extends WorkflowSummary {
  definition: WorkflowDefinition;
}

export interface WorkflowDefinition {
  initial: string;
  states: {
    [key: string]: State;
  };
}

export interface State {
  label: string;
  description: string;
  events: {
    [key: string]: Event;
  };
  timer?: Timer;
}

export interface Timer {
  duration: string;
  on_expire: string;
}

export interface Event {
  label: string;
  next_state: string;
}

export interface ActiveWorkflow {
  id: number;
  workflow: number;
  date_submitted: string;

  current_state: string;
  current_deadline?: string;

  definition: WorkflowDefinition;
}

export interface WorkflowListRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;
}

export type WorkflowListResponse = {
  // The total number of workflows returned without a limit or offset
  total: number;
  workflows: WorkflowSummary[];
};

export type WorkflowDetailsResponse = Workflow;
