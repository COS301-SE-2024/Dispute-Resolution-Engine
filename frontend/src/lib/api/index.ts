import { Result } from "../types";
import { API_URL } from "../utils";

export type Country = {
  id: number;
  country_code: string;
  country_name: string;
};

export async function fetchCountries(): Promise<Result<Country[]>> {
  const result = await fetch(`${API_URL}/utils/countries`, {
    cache: "no-store",
  });
  return await result.json();
}
