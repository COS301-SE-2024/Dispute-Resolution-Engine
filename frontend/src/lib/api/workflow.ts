"use server";

import { WorkflowListRequest, WorkflowListResponse } from "../interfaces/workflow";
import { sf, validateResult } from "../util";
import { getAuthToken } from "../util/jwt";
import { API_URL } from "../utils";

export async function getWorkflowList(req: WorkflowListRequest): Promise<WorkflowListResponse> {
  return sf(`${API_URL}/workflows`, {
    method: "POST",
    body: JSON.stringify(req),
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<WorkflowListResponse>);
}
