"use server";

import { Result } from "@/lib/types";

const API_URL = process.env.API_URL;

export type DisputeSummary = {
  id: string;
  title: string;
};

export async function fetchDisputes(user: string): Promise<Result<DisputeSummary[]>> {
  const response: Result<DisputeSummary[]> = await fetch(`${API_URL}/api`, {
    cache: "no-store",
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
