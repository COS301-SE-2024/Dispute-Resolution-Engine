"use server";

import { getAuthToken } from "../jwt";
import { ObjectionStatus, type ObjectionListResponse } from "../types/experts";
import { API_URL, sf, validateResult } from "../utils";

export async function getExpertObjections(dispute: number): Promise<ObjectionListResponse> {
  return sf(`${API_URL}/disputes/experts/objections`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      dispute_id: dispute,
    }),
  }).then(validateResult<ObjectionListResponse>);
}

export async function changeObjectionStatus(
  objection: number,
  status: ObjectionStatus
): Promise<void> {
  await sf(`${API_URL}/disputes/objections/${objection}`, {
    method: "PATCH",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      status,
    }),
  });
}
