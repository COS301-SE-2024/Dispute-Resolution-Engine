"use server";

import { getAuthToken } from "../jwt";
import {
  type TicketDetailsResponse,
  type TicketListRequest,
  type TicketListResponse,
  type TicketMessage,
  type TicketMessageResponse,
  type TicketStatus,
  type Ticket,
} from "../types/tickets";
import { API_URL, sf, validateResult } from "../utils";

export async function getTicketSummaries(req: TicketListRequest): Promise<TicketListResponse> {
  console.log(req);
  return sf(`${API_URL}/tickets`, {
    method: "POST",
    body: JSON.stringify(req),
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<TicketListResponse>);
}

export async function getTicketDetails(id: number): Promise<Ticket> {
  return sf(`${API_URL}/tickets/${id}`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  })
    .then(validateResult<TicketDetailsResponse>)
    .then(async (res) => {
      console.log(res);
      return res;
    });
}

export async function changeTicketStatus(id: string, status: TicketStatus): Promise<void> {
  await sf(`${API_URL}/tickets/${id}`, {
    method: "PATCH",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      status,
    }),
  });
}

export async function addTicketMessage(id: string, message: string): Promise<TicketMessage> {
  return sf(`${API_URL}/tickets/${id}/messages`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify({
      message,
    }),
  }).then(validateResult<TicketMessageResponse>);
}
