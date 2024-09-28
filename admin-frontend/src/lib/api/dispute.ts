"use server";

import {
  ActiveWorkflow,
  AdminDisputesRequest,
  AdminDisputesResponse,
  DisputeDetailsResponse,
  DisputeStatus,
} from "@/lib/types";
import { API_URL, sf, validateResult } from "../utils";
import { getAuthToken } from "../jwt";

export async function getDisputeList(req: AdminDisputesRequest): Promise<AdminDisputesResponse> {
  return sf(`${API_URL}/disputes`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      ...req,
      sort: {
        attr: "title",
      },
    }),
  }).then(validateResult<AdminDisputesResponse>);
}

export async function getDisputeDetails(id: number): Promise<DisputeDetailsResponse> {
  return sf(`${API_URL}/disputes/${id}`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  })
    .then(validateResult<DisputeDetailsResponse>)
    .then((res) => {
      console.log(res);
      return res;
    });
}

export async function changeDisputeStatus(id: number, status: DisputeStatus): Promise<void> {
  await sf(`${API_URL}/disputes/${id}/status`, {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      status,
    }),
  });
}

export async function getDisputeWorkflow(id: number): Promise<ActiveWorkflow> {
  return sf(`${API_URL}/disputes/${id}/workflow`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<ActiveWorkflow>);
}

export async function changeDisputeState(id: number, state: string): Promise<void> {
  await sf(`${API_URL}/workflows/reset`, {
    method: "PATCH",
    body: JSON.stringify({
      dispute_id: id,
      state,
    }),
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  });
}
