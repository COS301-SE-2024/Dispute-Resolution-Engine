"use server";

import { Result } from "@/lib/types";
import { API_URL } from "@/lib/utils";

export type DisputeSummary = {
  id: string;
  title: string;
};

export async function fetchDisputes(user: string): Promise<Result<DisputeSummary[]>> {
  const response: Result<DisputeSummary[]> = await fetch(`${API_URL}/api`, {
    cache: "no-cache",
    method: "POST",
    body: JSON.stringify({
      request_type: "dispute_summary",
      body: {
        id: user,
      },
    }),
  }).then((res) => res.json());
  return response;
}
