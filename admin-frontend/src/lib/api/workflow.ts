"use server";

import { ReactFlowJsonObject } from "@xyflow/react";

import { Workflow, GraphState, GraphTrigger } from "@/lib/types";

import {
  type WorkflowDetailsResponse,
  type WorkflowUpdateRequest,
  type WorkflowCreateRequest,
  type WorkflowCreateResponse,
  type WorkflowListRequest,
  type WorkflowListResponse,
  WorkflowDefinition,
} from "../types/workflow";

import { getAuthToken } from "../jwt";
import { API_URL, sf, validateResult } from "../utils";

export async function createWorkflow(req: WorkflowCreateRequest): Promise<WorkflowCreateResponse> {
  console.log("Request in createWorkflow", JSON.stringify(req));
  return sf(`${API_URL}/workflows/create`, {
    method: "POST",
    body: JSON.stringify(req),
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<WorkflowCreateResponse>);
}

export async function getWorkflowList(req: WorkflowListRequest): Promise<WorkflowListResponse> {
  return sf(`${API_URL}/workflows`, {
    method: "POST",
    body: JSON.stringify(req),
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  })
    .then(validateResult<WorkflowListResponse>)

    .then((res) => {
      console.log(res);
      return res;
    });
}

export async function getWorkflowDetails(id: number): Promise<WorkflowDetailsResponse> {
  return sf(`${API_URL}/workflows/${id}`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<WorkflowDetailsResponse>);
}

export async function updateWorkflow(id: number, req: WorkflowUpdateRequest): Promise<void> {
  await sf(`${API_URL}/workflows/${id}`, {
    method: "PATCH",
    body: JSON.stringify(req),
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  });
}

export async function deleteWorkflow(id: string): Promise<void> {
  await sf(`${API_URL}/workflows/${id}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  });
}

export async function graphToWorkflow({
  nodes,
  edges,
}: ReactFlowJsonObject<GraphState, GraphTrigger>): Promise<WorkflowDefinition> {
  // console.log(nodes, edges);
  const initial = nodes.find((state) => state.data.initial)?.id;
  if (!initial) {
    throw new Error("No initial state!");
  }

  return {
    initial: initial,
    states: Object.fromEntries(
      nodes.map((node) => [
        node.id,
        {
          label: node.data.label,
          description: node.data.description,
          events: Object.fromEntries(
            edges
              .filter((edge) => edge.source == node.id)
              .map((edge) => [
                edge.id,
                {
                  label: edge.data ? edge.data.trigger : "default",
                  next_state: edge.target,
                },
              ])
          ),
        },
      ])
    ),
  };
}

export async function workflowToGraph(
  workflow: WorkflowDefinition
): Promise<[GraphState[], GraphTrigger[]]> {
  let triggers: GraphTrigger[] = [];
  let currId: number = 0;
  for (const stateKey in workflow.states) {
    const state = workflow.states[stateKey];

    for (const eventKey in state.events) {
      const event = state.events[eventKey];
      const trigger: GraphTrigger = {
        id: (currId++).toString(),
        source: stateKey,
        target: event.next_state,
        data: { trigger: event.label },
        type: "custom-edge",
      };
      triggers.push(trigger);
    }
  }
  return [
    Object.keys(workflow.states).map((id) => ({
      id,
      type: "customNode",
      data: {
        label: workflow.states[id].label,
        initial: workflow.initial == id ? true : undefined,
        description: workflow.states[id].description,
        edges: [],
      },
      position: { x: 0, y: 0 },
    })),
    triggers,
  ];
}
