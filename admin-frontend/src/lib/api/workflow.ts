"use server";

import {
  type WorkflowDetailsResponse,
  type WorkflowUpdateRequest,
  type WorkflowCreateRequest,
  type WorkflowCreateResponse,
  type WorkflowListRequest,
  type WorkflowListResponse,
  Workflow,
} from "../types/workflow";

// import { getAuthToken } from "../jwt";
// import { API_URL, sf, validateResult } from "../utils";

let MOCK_WORKFLOWS: Workflow[] = [
  {
    id: 1,
    name: "User Onboarding",
    date_created: "2023-08-01",
    last_updated: "2023-09-01",
    author: {
      id: 101,
      full_name: "John Doe",
    },
    definition: {
      initial: "start",
      states: {
        start: {
          label: "Start",
          description: "Initial state of onboarding",
          events: {
            next: {
              label: "Proceed to Details",
              next_state: "details",
            },
          },
        },
        details: {
          label: "Enter Details",
          description: "User enters their details",
          events: {
            submit: {
              label: "Submit Details",
              next_state: "verification",
            },
          },
        },
        verification: {
          label: "Verify",
          description: "Verification of submitted details",
          events: {
            approve: {
              label: "Approve Details",
              next_state: "completed",
            },
            reject: {
              label: "Reject Details",
              next_state: "start",
            },
          },
          timer: {
            duration: "24h",
            on_expire: "start",
          },
        },
        completed: {
          label: "Completed",
          description: "Onboarding process completed",
          events: {},
        },
      },
    },
  },
  {
    id: 2,
    name: "Order Fulfillment",
    date_created: "2023-07-15",
    last_updated: "2023-08-01",
    author: {
      id: 102,
      full_name: "Jane Smith",
    },
    definition: {
      initial: "pending",
      states: {
        pending: {
          label: "Pending",
          description: "Order is awaiting fulfillment",
          events: {
            process: {
              label: "Start Processing",
              next_state: "processing",
            },
          },
        },
        processing: {
          label: "Processing",
          description: "Order is being processed",
          events: {
            ship: {
              label: "Ship Order",
              next_state: "shipped",
            },
          },
        },
        shipped: {
          label: "Shipped",
          description: "Order has been shipped",
          events: {
            deliver: {
              label: "Deliver Order",
              next_state: "delivered",
            },
          },
        },
        delivered: {
          label: "Delivered",
          description: "Order has been delivered",
          events: {},
        },
      },
    },
  },
  {
    id: 3,
    name: "Incident Management",
    date_created: "2023-06-20",
    last_updated: "2023-07-10",
    author: {
      id: 103,
      full_name: "Alice Johnson",
    },
    definition: {
      initial: "reported",
      states: {
        reported: {
          label: "Reported",
          description: "Incident has been reported",
          events: {
            assign: {
              label: "Assign to Team",
              next_state: "assigned",
            },
          },
        },
        assigned: {
          label: "Assigned",
          description: "Incident assigned to a team",
          events: {
            resolve: {
              label: "Resolve Incident",
              next_state: "resolved",
            },
            escalate: {
              label: "Escalate Incident",
              next_state: "escalated",
            },
          },
        },
        escalated: {
          label: "Escalated",
          description: "Incident has been escalated",
          events: {
            resolve: {
              label: "Resolve Escalated Incident",
              next_state: "resolved",
            },
          },
          timer: {
            duration: "48h",
            on_expire: "assigned",
          },
        },
        resolved: {
          label: "Resolved",
          description: "Incident has been resolved",
          events: {},
        },
      },
    },
  },
  {
    id: 4,
    name: "Content Approval",
    date_created: "2023-07-01",
    last_updated: "2023-07-20",
    author: {
      id: 104,
      full_name: "Bob Williams",
    },
    definition: {
      initial: "draft",
      states: {
        draft: {
          label: "Draft",
          description: "Content is in draft mode",
          events: {
            submit: {
              label: "Submit for Review",
              next_state: "review",
            },
          },
        },
        review: {
          label: "Review",
          description: "Content is under review",
          events: {
            approve: {
              label: "Approve Content",
              next_state: "published",
            },
            reject: {
              label: "Reject Content",
              next_state: "draft",
            },
          },
          timer: {
            duration: "48h",
            on_expire: "draft",
          },
        },
        published: {
          label: "Published",
          description: "Content is published",
          events: {},
        },
      },
    },
  },
  {
    id: 5,
    name: "Employee Offboarding",
    date_created: "2023-08-05",
    last_updated: "2023-09-01",
    author: {
      id: 105,
      full_name: "Carol Brown",
    },
    definition: {
      initial: "notice_given",
      states: {
        notice_given: {
          label: "Notice Given",
          description: "Employee has given notice",
          events: {
            process: {
              label: "Start Offboarding",
              next_state: "offboarding",
            },
          },
        },
        offboarding: {
          label: "Offboarding",
          description: "Offboarding process is ongoing",
          events: {
            complete: {
              label: "Complete Offboarding",
              next_state: "completed",
            },
          },
        },
        completed: {
          label: "Completed",
          description: "Offboarding process is complete",
          events: {},
        },
      },
    },
  },
  {
    id: 6,
    name: "Customer Support Ticket",
    date_created: "2023-08-10",
    last_updated: "2023-09-05",
    author: {
      id: 106,
      full_name: "David Lee",
    },
    definition: {
      initial: "open",
      states: {
        open: {
          label: "Open",
          description: "Support ticket is opened",
          events: {
            assign: {
              label: "Assign Ticket",
              next_state: "assigned",
            },
          },
        },
        assigned: {
          label: "Assigned",
          description: "Ticket is assigned to an agent",
          events: {
            resolve: {
              label: "Resolve Ticket",
              next_state: "resolved",
            },
          },
        },
        resolved: {
          label: "Resolved",
          description: "Ticket has been resolved",
          events: {
            reopen: {
              label: "Reopen Ticket",
              next_state: "open",
            },
          },
        },
      },
    },
  },
  {
    id: 7,
    name: "Purchase Order Approval",
    date_created: "2023-08-15",
    last_updated: "2023-09-10",
    author: {
      id: 107,
      full_name: "Emily White",
    },
    definition: {
      initial: "created",
      states: {
        created: {
          label: "Created",
          description: "Purchase order created",
          events: {
            submit: {
              label: "Submit for Approval",
              next_state: "approval",
            },
          },
        },
        approval: {
          label: "Approval",
          description: "Waiting for manager approval",
          events: {
            approve: {
              label: "Approve",
              next_state: "approved",
            },
            reject: {
              label: "Reject",
              next_state: "created",
            },
          },
          timer: {
            duration: "72h",
            on_expire: "created",
          },
        },
        approved: {
          label: "Approved",
          description: "Purchase order approved",
          events: {},
        },
      },
    },
  },
  {
    id: 8,
    name: "Project Management",
    date_created: "2023-09-01",
    last_updated: "2023-09-15",
    author: {
      id: 108,
      full_name: "Frank Miller",
    },
    definition: {
      initial: "initiated",
      states: {
        initiated: {
          label: "Initiated",
          description: "Project is initiated",
          events: {
            plan: {
              label: "Start Planning",
              next_state: "planning",
            },
          },
        },
        planning: {
          label: "Planning",
          description: "Project planning phase",
          events: {
            execute: {
              label: "Execute Project",
              next_state: "execution",
            },
          },
        },
        execution: {
          label: "Execution",
          description: "Project execution phase",
          events: {
            complete: {
              label: "Complete Project",
              next_state: "completed",
            },
          },
        },
        completed: {
          label: "Completed",
          description: "Project is completed",
          events: {},
        },
      },
    },
  },
  {
    id: 9,
    name: "Feature Development",
    date_created: "2023-07-25",
    last_updated: "2023-08-10",
    author: {
      id: 109,
      full_name: "Grace Young",
    },
    definition: {
      initial: "idea",
      states: {
        idea: {
          label: "Idea",
          description: "Feature idea generation",
          events: {
            plan: {
              label: "Plan Feature",
              next_state: "planning",
            },
          },
        },
        planning: {
          label: "Planning",
          description: "Feature planning phase",
          events: {
            develop: {
              label: "Start Development",
              next_state: "development",
            },
          },
        },
        development: {
          label: "Development",
          description: "Feature is under development",
          events: {
            review: {
              label: "Review Feature",
              next_state: "review",
            },
          },
        },
        review: {
          label: "Review",
          description: "Feature is being reviewed",
          events: {
            approve: {
              label: "Approve Feature",
              next_state: "completed",
            },
            reject: {
              label: "Reject Feature",
              next_state: "development",
            },
          },
        },
        completed: {
          label: "Completed",
          description: "Feature development is complete",
          events: {},
        },
      },
    },
  },
  {
    id: 10,
    name: "Expense Reimbursement",
    date_created: "2023-06-30",
    last_updated: "2023-07-25",
    author: {
      id: 110,
      full_name: "Henry Green",
    },
    definition: {
      initial: "submitted",
      states: {
        submitted: {
          label: "Submitted",
          description: "Reimbursement request submitted",
          events: {
            review: {
              label: "Review Request",
              next_state: "review",
            },
          },
        },
        review: {
          label: "Review",
          description: "Reimbursement request is under review",
          events: {
            approve: {
              label: "Approve Request",
              next_state: "approved",
            },
            reject: {
              label: "Reject Request",
              next_state: "submitted",
            },
          },
          timer: {
            duration: "14d",
            on_expire: "submitted",
          },
        },
        approved: {
          label: "Approved",
          description: "Reimbursement has been approved",
          events: {},
        },
      },
    },
  },
];

export async function createWorkflow(req: WorkflowCreateRequest): Promise<WorkflowCreateResponse> {
  // return sf(`${API_URL}/workflows`, {
  //   method: "POST",
  //   body: JSON.stringify(req),
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // }).then(validateResult<WorkflowCreateResponse>);
  throw new Error("not implemented");
}

export async function getWorkflowList(req: WorkflowListRequest): Promise<WorkflowListResponse> {
  // return sf(`${API_URL}/workflows`, {
  //   method: "POST",
  //   body: JSON.stringify(req),
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // }).then(validateResult<WorkflowListResponse>);

  let data = MOCK_WORKFLOWS;
  let count = data.length;
  if (req.search) {
    data = data.filter((wf) =>
      [wf.name, wf.author.full_name].join(" ").toLowerCase().includes(req.search!.toLowerCase())
    );
  }

  data = data.slice(req.offset ?? 0, req.limit ? (req.offset ?? 0) + req.limit : undefined);

  return {
    total: count,
    workflows: data,
  };
}

export async function getWorkflowDetails(id: number): Promise<WorkflowDetailsResponse> {
  // return sf(`${API_URL}/workflows/${id}`, {
  //   method: "GET",
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // }).then(validateResult<WorkflowDetailsResponse>);

  const workflow = MOCK_WORKFLOWS.find((wf) => wf.id === id);
  if (!workflow) {
    throw new Error("Workflow not found");
  }
  return workflow;
}

export async function updateWorkflow(id: string, req: WorkflowUpdateRequest): Promise<void> {
  // await sf(`${API_URL}/workflows/${id}`, {
  //   method: "PATCH",
  //   body: JSON.stringify(req),
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // });
  if (MOCK_WORKFLOWS.every((wf) => wf.id === parseInt(id))) {
    throw new Error("Workflow not found");
  }
  MOCK_WORKFLOWS.filter((wf) =>
    wf.id === parseInt(id)
      ? {
          ...wf,
          ...req,
        }
      : wf
  );
}

export async function deleteWorkflow(id: string): Promise<void> {
  // await sf(`${API_URL}/workflows/${id}`, {
  //   method: "DELETE",
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // });
  if (MOCK_WORKFLOWS.every((wf) => wf.id === parseInt(id))) {
    throw new Error("Workflow not found");
  }
  MOCK_WORKFLOWS = MOCK_WORKFLOWS.filter((wf) => wf.id !== parseInt(id));
}
