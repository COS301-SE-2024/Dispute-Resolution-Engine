import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { Result, resultSchema } from "./types";

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
export function sf<T>(input: string | URL | globalThis.Request, init?: RequestInit): Promise<T> {
  return fetch(input, init)
    .then((res) => res.json())
    .then((res) => {
      const { data, error } = resultSchema.safeParse(res);
      if (error) {
        console.log(error);
        throw new Error(error.issues[0].message);
      }
      if (data.error) {
        throw new Error(data.error);
      }
      return data.data! as T;
    });
}
export function resultify<T>(prom: Promise<T>): Promise<Result<T>> {
  return prom
    .then((data) => ({ data }))
    .catch((e: Error) => ({
      error: e.message,
    }));
}

export const API_URL = process.env.API_URL;
