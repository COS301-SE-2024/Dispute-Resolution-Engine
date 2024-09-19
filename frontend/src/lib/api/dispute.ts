"use server";

import { Result } from "@/lib/types";
import {
  DisputeListResponse,
  DisputeResponse,
  DisputeStatusUpdateRequest,
} from "../interfaces/dispute";
import { cookies } from "next/headers";
import { JWT_KEY } from "../constants";
import { API_URL } from "@/lib/utils";
import { getAuthToken } from "../util/jwt";
import { revalidatePath } from "next/cache";

export async function getDisputeList(): Promise<Result<DisputeListResponse>> {
  const res = await fetch(`${API_URL}/disputes`, {
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  })
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
  return res;
}

export async function getDisputeDetails(id: string): Promise<Result<DisputeResponse>> {
  const res = await fetch(`${API_URL}/disputes/${id}`, {
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  })
    .then(function (res) {
      return res.json();
    })
    .catch((e: Error) => ({
      error: e.message,
    }));

  console.log(res);
  return res;
}
export async function updateDisputeStatus(
  id: string,
  status: string,
): Promise<Result<DisputeResponse>> {
  const body: DisputeStatusUpdateRequest = { dispute_id: id, status };
  const res = await fetch(`${API_URL}/disputes/dispute/status`, {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify(body),
  })
    .then(function (res) {
      return res.json();
    })
    .catch((e: Error) => ({
      error: e.message,
    }));
  console.log("RESPONSE IN UPDATE DISPUTE\n", res);
  console.log("BODY WAS\n", JSON.stringify(body));
  revalidatePath(`/disputes/${id}`);
  return res;
}
export async function getStatusEnum(): Promise<string[]> {
  const res = await fetch(`${API_URL}/utils/dispute_statuses`, {
    method: "GET",
  })
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
  return res.data;
}
