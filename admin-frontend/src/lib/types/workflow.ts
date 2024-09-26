export interface WorkflowSummary {
  id: number;
  name: string;
  date_created: string;
  last_updated: string;
  author: {
    id: number;
    full_name: string;
  };
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

export interface ActiveWorkflow extends Workflow {
  current_state: {
    // The ID of the current state
    id: string;

    // The deadline of the current state (if any)
    deadline?: string;
  };
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

export interface WorkflowUpdateRequest {
  // The new name of the workflow (or unchanged)
  name?: string;

  // Updates the tags assigned to the workflow (will overwrite existing tags)
  // tags?: string[];

  // Updates the workflow definition
  definition?: WorkflowDefinition;
}

export interface WorkflowCreateRequest {
  // The name of the workflow
  name: string;

  // The tags assigned to the workflow
  tags: string[];

  // The workflow definition
  definition: WorkflowDefinition;
}

export type WorkflowCreateResponse = Workflow;

export type ActiveWorkflowResponse = ActiveWorkflow;