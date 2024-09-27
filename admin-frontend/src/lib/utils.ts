import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { Result, resultSchema } from "./types";

export const API_URL = process.env.API_URL;

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function unwrapResult<T>(p: Promise<Result<T>>): Promise<T> {
  return p.then(({ data, error }) => {
    if (error) {
      throw new Error(error);
    }
    return data!;
  });
}

/**
 * Wrapper around fetch that automatically parses the return result and throws an error if something went wrong
 */
export function resultify<T>(prom: Promise<T>): Promise<Result<T>> {
  return prom
    .then((data) => ({ data }))
    .catch((e: Error) => ({
      error: e.message,
    }));
}



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
