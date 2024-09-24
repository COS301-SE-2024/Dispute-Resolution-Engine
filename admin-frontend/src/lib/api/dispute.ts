"use server";

import { Result } from "@/lib/types";
import {
  AdminDisputesRequest,
  AdminDisputesResponse,
  DisputeDetailsResponse,
  DisputeStatus,
} from "@/lib/types";
import { API_URL, resultify, sf } from "../utils";
import { getAuthToken } from "../jwt";

export async function getDisputeList(
  req: AdminDisputesRequest
): Promise<Result<AdminDisputesResponse>> {
  const res = await resultify(
    sf<AdminDisputesResponse>(`${API_URL}/disputes`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify({
        ...req,
      }),
    })
  );
  console.log(res);
  return res;
}

export async function getDisputeDetails(id: string): Promise<Result<DisputeDetailsResponse>> {
  const res = await resultify(
    sf<DisputeDetailsResponse>(`${API_URL}/disputes/${id}`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
    })
  );
  return res;
}

export async function changeDisputeStatus(
  id: string,
  status: DisputeStatus
): Promise<Result<string>> {
  const res = await resultify(
    sf<string>(`${API_URL}/disputes/${id}/status`, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify({
        status,
      }),
    })
  );
  return res;
}
export async function deleteEvidence(
  disputeId: string,
  evidenceId: string
): Promise<Result<string>> {
  return delayResolve(
    {
      data: "Evidence deleted",
    },
    1000
  );
}

function delayResolve<T>(data: T, millis: number): Promise<T> {
  return new Promise((res) => {
    setTimeout(() => res(data), millis);
  });
}
