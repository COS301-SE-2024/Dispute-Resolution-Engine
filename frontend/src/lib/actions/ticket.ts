"use server";

import { TicketCreateResponse } from "../interfaces/ticket";
import { CreateTicketData, CreateTicketErrors, createTicketSchema } from "../schema/ticket";
import { Result } from "../types";
import { getAuthToken } from "../util/jwt";
import { API_URL, formFetch } from "../utils";

export async function createTicket(
  _initial: unknown,
  data: FormData
): Promise<Result<string, CreateTicketErrors>> {
  const { data: parsed, error: parseErr } = createTicketSchema.safeParse(Object.fromEntries(data));
  if (parseErr) {
    return {
      error: parseErr.format(),
    };
  }

  const res = await formFetch<CreateTicketData, TicketCreateResponse>(`${API_URL}/tickets/create`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: JSON.stringify(parsed),
  });
  if (res.error) {
    return {
      error: res.error,
    };
  }
  return {
    data: "TODO",
  };
}
