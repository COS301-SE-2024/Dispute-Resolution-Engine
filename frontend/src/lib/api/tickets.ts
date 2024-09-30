"use server";

import { revalidatePath } from "next/cache";
import {
  Ticket,
  TicketDetailsResponse,
  TicketListResponse,
  TicketMessage,
  TicketMessageResponse,
} from "../interfaces/tickets";
import { sf, validateResult } from "../util";
import { getAuthToken } from "../util/jwt";
import { API_URL } from "../utils";

export async function getTicketSummaries(dispute: number): Promise<TicketListResponse> {
  return sf(`${API_URL}/tickets`, {
    method: "POST",
    body: JSON.stringify({
      // filter: [
      //   {
      //     attr: "dispute_id",
      //     value: dispute.toString(),
      //   },
      // ],
    }),
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
  }).then(validateResult<TicketDetailsResponse>);
}

export async function addTicketMessage(id: number, message: string): Promise<TicketMessage> {
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
