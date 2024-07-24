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

export async function getDisputeList(): Promise<Result<DisputeListResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }

  const res = await fetch(`${API_URL}/disputes`, {
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
  })
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
  return res;
}

export async function getDisputeDetails(id: string): Promise<Result<DisputeResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }

  const res = await fetch(`${API_URL}/disputes/${id}`, {
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
  })
    .then(function (res) {
      return res.json();
    })
    .catch((e: Error) => ({
      error: e.message,
    }));
  return res;
}
export async function updateDisputeStatus(
  id: string,
  status: string
): Promise<Result<DisputeResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }
  const body: DisputeStatusUpdateRequest = { id, status };
  return fetch(`${API_URL}/dispute/status`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
    body: JSON.stringify(body),
  })
    .then(function (res) {
      return res.json();
    })
    .catch((e: Error) => ({
      error: e.message,
    }));
}
