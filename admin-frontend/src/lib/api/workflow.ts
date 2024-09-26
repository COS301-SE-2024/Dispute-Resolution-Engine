"use server";

import {
  type WorkflowDetailsResponse,
  type WorkflowUpdateRequest,
  type WorkflowCreateRequest,
  type WorkflowCreateResponse,
  type WorkflowListRequest,
  type WorkflowListResponse,
} from "../types/workflow";

// import { getAuthToken } from "../jwt";
// import { API_URL, sf, validateResult } from "../utils";

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
  throw new Error("not implemented");
}

export async function getWorkflowDetails(id: string): Promise<WorkflowDetailsResponse> {
  // return sf(`${API_URL}/workflows/${id}`, {
  //   method: "GET",
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // }).then(validateResult<WorkflowDetailsResponse>);
  throw new Error("not implemented");
}

export async function updateWorkflow(id: string, req: WorkflowUpdateRequest): Promise<void> {
  // await sf(`${API_URL}/workflows/${id}`, {
  //   method: "PATCH",
  //   body: JSON.stringify(req),
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // });
  throw new Error("not implemented");
}

export async function deleteWorkflow(id: string): Promise<void> {
  // await sf(`${API_URL}/workflows/${id}`, {
  //   method: "DELETE",
  //   headers: {
  //     Authorization: `Bearer ${getAuthToken()}`,
  //   },
  // });
  throw new Error("not implemented");
}
