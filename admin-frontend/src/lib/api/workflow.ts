"use server";

import {
  type WorkflowDetailsResponse,
  type WorkflowUpdateRequest,
  type WorkflowCreateRequest,
  type WorkflowCreateResponse,
  type WorkflowListRequest,
  type WorkflowListResponse,
} from "../types/workflow";

export async function createWorkflow(req: WorkflowCreateRequest): Promise<WorkflowCreateResponse> {
  throw new Error("not implemented");
}

export async function getWorkflowList(req: WorkflowListRequest): Promise<WorkflowListResponse> {
  throw new Error("not implemented");
}

export async function getWorkflowDetails(id: string): Promise<WorkflowDetailsResponse> {
  throw new Error("not implemented");
}

export async function updateWorkflow(req: WorkflowUpdateRequest): Promise<void> {
  throw new Error("not implemented");
}

export async function deleteWorkflow(id: string): Promise<void> {
  throw new Error("not implemented");
}
