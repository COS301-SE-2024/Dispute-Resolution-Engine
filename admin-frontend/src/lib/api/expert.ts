"use server";

import { getAuthToken } from "../jwt";
import { ObjectionStatus, type ObjectionListResponse } from "../types/experts";
import { API_URL, sf, validateResult } from "../utils";

export async function getExpertObjections(dispute: number): Promise<ObjectionListResponse> {
  return sf(`${API_URL}/disputes/${dispute}/objections`, {
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<ObjectionListResponse>);
}

export async function changeObjectionStatus(
  dispute: number,
  objection: number,
  status: ObjectionStatus
): Promise<void> {
  await sf(`${API_URL}/disputes/${dispute}/objections/${objection}`, {
    method: "PATCH",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      status,
    }),
  });
}
