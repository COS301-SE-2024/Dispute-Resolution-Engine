import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { z } from "zod";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const API_URL = process.env.API_URL;

export const resultSchema = z.object({
  data: z.any().optional(),
  error: z.string().optional(),
});

/**
 * Parses the response body into JSON and checks that it conforms to the Result type
 */
export async function validateResult<T>(res: Response): Promise<T> {
  const json = await res.json();
  const { data, error } = resultSchema.safeParse(json);
  if (error) {
    console.log(error);
    throw new Error(error.issues[0].message);
  }
  if (data.error) {
    throw new Error(data.error);
  }
  return data.data! as T;
}

/**
 * Wrapper around fetch that checks if the request was fine
 */
export function sf(
  input: string | URL | globalThis.Request,
  init?: RequestInit
): Promise<Response> {
  return fetch(input, init).then((res) => {
    if (!res.ok) {
      throw new Error(`Request failed with code ${res.status}`);
    }
    return res;
  });
}
