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
  return {
    data: [],
  };

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
  return {
    data: {
      id: "1",
      title: "Mock title",
      description: "Description",
      status: "In porgress",
      case_date: "today",
      role: "Complainant",

      evidence: [],
      experts: [],
    },
  };

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
  status: string,
): Promise<Result<DisputeResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }
  const body: DisputeStatusUpdateRequest = { dispute_id: id, status };
  const res = await fetch(`${API_URL}/disputes/dispute/status`, {
    method: "PUT",
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
  console.log("RESPONSE IN UPDATE DISPUTE\n", res);
  console.log("BODY WAS\n", JSON.stringify(body));
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
