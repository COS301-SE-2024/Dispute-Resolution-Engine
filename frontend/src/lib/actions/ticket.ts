"use server";

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

  const res = await formFetch<CreateTicketData, string>(
    `${API_URL}/disputes/${parsed.dispute}/tickets`,
    {
      method: "POST",
      headers: {
        // Sub this for the proper getAuthToken thing
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify({
        subject: parsed.subject,
        body: parsed.body,
      }),
    }
  );
  if (res.error) {
    return {
      error: res.error,
    };
  }
  return {
    data: "Ticket created",
  };
}
