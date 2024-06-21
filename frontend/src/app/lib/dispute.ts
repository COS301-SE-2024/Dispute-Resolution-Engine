"use server";

import { Result } from "@/lib/types";
import { API_URL } from "@/lib/utils";

export type DisputeSummary = {
  id: string;
  title: string;
};

export async function fetchDisputes(user: string): Promise<Result<DisputeSummary[]>> {
  return {
    data: [...Array(10).keys()].map((i) => ({
      id: i.toString(),
      title: `Dispute #${i}`,
    })),
  };
  // const response: Result<DisputeSummary[]> = await fetch(`${API_URL}/api`, {
  //   cache: "no-store",
  //   method: "POST",
  //   body: JSON.stringify({
  //     request_type: "dispute_summary",
  //     body: {
  //       id: user,
  //     },
  //   }),
  // }).then((res) => res.json());
  // return response;
}
