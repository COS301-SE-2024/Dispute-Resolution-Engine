"use server";

import { Result } from "../types";
import { sf, validateResult } from "../util";
import { API_URL } from "../utils";

export type Country = {
  id: number;
  country_code: string;
  country_name: string;
};

export type Language = {
  id: string;
  label: string;
};

export async function getCountries(): Promise<Country[]> {
  return sf(`${API_URL}/utils/countries`, {
    cache: "no-store",
  }).then(validateResult<Country[]>);
}

export async function fetchLanguages(): Promise<Language[]> {
  return [
    {
      id: "en-US",
      label: "English",
    },
  ];
}
