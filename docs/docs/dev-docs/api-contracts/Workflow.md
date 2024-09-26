# Workflows

The following types are used for describing workflows:

```ts
interface WorkflowSummary {
  id: string;
  name: string;
  date_created: string;
  last_updated: string;
  author: {
    id: string;
    full_name: string;
  };
}

interface Workflow extends WorkflowSummary {
  definition: WorkflowDefinition;
}

interface ActiveWorkflow extends Workflow {
  current_state: {
    // The ID of the current state
    id: string;

    // The deadline of the current state (if any)
    deadline?: string;
  };
}

interface WorkflowDefinition {
  initial: string;
  states: {
    [key: string]: State;
  };
}

interface State {
  label: string;
  description: string;
  events: {
    [key: string]: Event;
  };
  timer?: Timer;
}

interface Timer {
  duration: string;
  on_expire: string;
}

interface Event {
  label: string;
  next_state: string;
}
```

## List of workflows

- **Endpoint:** `POST /workflows`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accessible by an administrator

```ts
interface WorkflowListRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;
}

type WorkflowListResponse = {
  // The total number of workflows returned without a limit or offset
  total: number;
  workflows: WorkflowSummary[];
};
```

## Workflow details

- **Endpoint:** `GET /workflows/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accessible by an administrator

```ts
type WorkflowDetailsResponse = Workflow;
```

## Updating a workflow

- **Endpoint:** `PATCH /workflows/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accessible by an administrator

```ts
interface WorkflowUpdateRequest {
  // The new name of the workflow (or unchanged)
  name?: string;

  // Updates the tags assigned to the workflow (will overwrite existing tags)
  tags?: string[];

  // Updates the workflow definition
  definition?: WorkflowDefinition;
}
```

On success, the resposne can jsut be an HTTP 204 (no content) response.

## Create a workflow

- **Endpoint:** `POST /workflows`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accessible by an administrator

```ts
interface WorkflowUpdateRequest {
  // The name of the workflow
  name: string;

  // The tags assigned to the workflow
  tags: string[];

  // The workflow definition
  definition: WorkflowDefinition;
}

type WorkflowUpdateResponse = Workflow;
```

## Deleting a workflow

- **Endpoint:** `DELETE /workflows/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accessible by an administrator

The endpoint does not accept any body. On success it will return a HTTP 204 (no content) message.

# Active Workflows

## Viewing all active workflows

## Viewing a particular workflow

Displays the full workflow details, along with the current state and timer

```ts
type ActiveWorkflowResponse = ActiveWorkflow;
```

## Patching an active workflow

Should be able to adjust the current state and optionally the timer
