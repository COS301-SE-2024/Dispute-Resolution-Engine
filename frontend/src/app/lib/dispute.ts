"use server";

const API_URL = process.env.API_URL;

export type DisputeSummary = {
    id: string;
    title: string;
}

export async function fetchDisputes(user: string): Promise<DisputeSummary[] | string> {
  const response = await fetch(`${API_URL}/api`, {
    method: "POST",
    body: JSON.stringify({
      request_type: "dispute_summary",
      body: {
        id: user,
      }
    })
  }).then(res => res.json());
  if (response.Error) {
    return response.Error;
  }
  return response;
}