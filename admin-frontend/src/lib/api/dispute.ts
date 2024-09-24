"use server";

import { Result } from "@/lib/types";
import {
  AdminDisputesRequest,
  AdminDisputesResponse,
  DisputeDetails,
  DisputeDetailsResponse,
  DisputeStatus,
} from "@/lib/types";
import { API_URL } from "../utils";
import { getAuthToken } from "../jwt";

export async function getDisputeList(
  req: AdminDisputesRequest
): Promise<Result<AdminDisputesResponse>> {
  const res = await fetch(`${API_URL}/disputes`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      ...req,
    }),
  })
    .then(async (res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
  console.log(res);
  return res;
}

export async function getDisputeDetails(id: string): Promise<Result<DisputeDetailsResponse>> {
  const res = await fetch(`${API_URL}/disputes/${id}`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  })
    .then(async (res) => {
      console.log(await res.clone().text());
      return res.json();
    })
    .catch((e: Error) => ({
      error: e.message,
    }));
  console.log(res);
  return res;
}

export async function changeDisputeStatus(
  id: string,
  status: DisputeStatus
): Promise<Result<string>> {
  return delayResolve(
    {
      data: "Status changed",
      // error: "Big bad error things",
    },
    1000
  );
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
