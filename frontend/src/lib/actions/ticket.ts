"use server";

import { CreateTicketErrors, createTicketSchema } from "../schema/ticket";
import { Result } from "../types";

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

  //   const res = await formFetch<CreateTicketData, string>(
  //     `${API_URL}/disputes/${parsed.dispute_id}/experts/reject`,
  //     {
  //       method: "POST",
  //       headers: {
  //         // Sub this for the proper getAuthToken thing
  //         Authorization: `Bearer ${getAuthToken()}`,
  //       },
  //       body: JSON.stringify({
  //         expert_id: parseInt(parsed.expert_id),
  //         reason: parsed.reason,
  //       }),
  //     }
  //   );

  //   if (!res.error) {
  //     revalidatePath(`/disputes/${parsed.dispute_id}`);
  //   }
  return {
    data: "cool",
  };
}
