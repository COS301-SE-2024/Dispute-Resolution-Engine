import { Result } from "../types";
import { API_URL } from "../utils";

export type Country = {
  id: number;
  country_code: string;
  country_name: string;
};

export async function fetchCountries(): Promise<Result<Country[]>> {
  return {
    data: [
      {
        id: 1,
        country_code: "za",
        country_name: "South Africa",
      },
    ],
  };
  // const result = await fetch(`${API_URL}/utils/countries`, {
  //   cache: "no-store",
  // });
  // console.log(await result.text());
  // return await result.json();
}
