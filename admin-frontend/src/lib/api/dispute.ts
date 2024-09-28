"use server";

import {
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

export async function getDisputeDetails(id: string): Promise<DisputeDetailsResponse> {
  return sf(`${API_URL}/disputes/${id}`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<DisputeDetailsResponse>);
}

export async function changeDisputeStatus(id: string, status: DisputeStatus): Promise<void> {
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
