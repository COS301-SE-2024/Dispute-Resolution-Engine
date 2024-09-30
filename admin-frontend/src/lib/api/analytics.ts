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
