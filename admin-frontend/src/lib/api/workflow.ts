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
  return sf(`${API_URL}/workflows`, {
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

export async function updateWorkflow(id: string, req: WorkflowUpdateRequest): Promise<void> {
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

export function graphToWorkflow({
  nodes,
  edges,
}: ReactFlowJsonObject<GraphState, GraphTrigger>): WorkflowDefinition {
  console.log(nodes, edges);
  return {
    initial: "Im not sure",
    states: Object.fromEntries(
      nodes.map((node) => [
        node.id,
        {
          label: node.data.label,
          description: "sure bud",
          events: Object.fromEntries(
            edges
              .filter((edge) => edge.source == node.id)
              .map((edge) => [
                edge.id,
                {
                  label: "oi blud, do somfin",
                  next_state: edge.target,
                },
              ])
          ),
        },
      ])
    ),
  };
}

export function workflowToGraph(workflow: Workflow): [GraphState[], GraphTrigger[]] {
  return [
    Object.keys(workflow.states).map((id) => ({
      id,
      type: "customNode",
      data: {
        label: workflow.states[id].label,
        edges: [],
      },
      position: { x: 0, y: 0 },
    })),
    [],
  ];
}
