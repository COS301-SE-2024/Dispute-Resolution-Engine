"use server";

import { getAuthToken } from "../jwt";
import { DisputeStatus } from "../types";
import { TicketStatus } from "../types/tickets";
import { API_URL, sf, validateResult } from "../utils";

export async function getDisputeCountByStatus(): Promise<Record<DisputeStatus, number>> {
  return sf(`${API_URL}/analytics/stats/disputes`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      group: "status",
    }),
  }).then(validateResult<Record<DisputeStatus, number>>);
}

export async function getTicketCountByStatus(): Promise<Record<TicketStatus, number>> {
  return sf(`${API_URL}/analytics/stats/tickets`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      group: "status",
    }),
  }).then(validateResult<Record<TicketStatus, number>>);
}

export async function getExpertsObjectionSummary(): Promise<Record<string, number>> {
  return sf(`${API_URL}/analytics/stats/expert_objections_view`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      group: "expert_full_name",
    }),
  }).then(validateResult<Record<string, number>>);
}

export async function getMonthlyDisputes(): Promise<Record<string, number>> {
  return sf(`${API_URL}/analytics/monthly/disputes`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      group: "case_date",
    }),
  })
    .then(validateResult<Record<string, number>>)
    .then((res) => {
      console.log(res);
      return res;
    });
}
