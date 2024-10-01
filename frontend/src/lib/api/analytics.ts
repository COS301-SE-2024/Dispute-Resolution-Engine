import { Result } from "../types";
import { sf, validateResult } from "../util";
import { getAuthToken } from "../util/jwt";
import { API_URL } from "../utils";

export type CountryCount = {
  nationality: string;
  count: number;
};

export async function fetchUserCountryDistribution(): Promise<Result<CountryCount[]>> {
  const jwt = getAuthToken();
  const res = await fetch(`${API_URL}/user/analytics`, {
    method: "POST",
    cache: "no-cache",
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
    body: JSON.stringify({
      groupBy: ["nationality"],
      count: true,
    }),
  });
  return res.json().catch(async (e: Error) => ({
    error: e.message,
  }));
}

export interface DisputeEstimate {
  days: number;
  hours: number;
  minutes: number;
}
export async function getDisputeEstimate(id: string): Promise<DisputeEstimate> {
  return sf(`${API_URL}/analytics/time/estimation`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  }).then(validateResult<DisputeEstimate>);
}
