import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { ZodFormattedError } from "zod";
import { Result } from "./types";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export type Slottable<P> = P & {
  asChild?: boolean;
};

export const API_URL = process.env.API_URL;

export async function resultify<T>(p: () => Promise<T>): Promise<Result<T>> {
  return p()
    .then((data) => ({ data }))
    .catch((err) => ({
      error: (err as Error).message,
    }));
}

export async function safeFetch<T>(
  input: string | URL | globalThis.Request,
  init?: RequestInit
): Promise<Result<T>> {
  return fetch(input, init)
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
}

export async function formFetch<TFormData, R = unknown>(
  input: string | URL | globalThis.Request,
  init?: RequestInit
): Promise<Result<R, ZodFormattedError<TFormData>>> {
  return fetch(input, init)
    .then((res) => res.json())
    .then((res) =>
      !res.error
        ? res
        : {
            error: {
              _errors: [res.error],
            },
          }
    )
    .catch((e: Error) => ({
      error: {
        _errors: [e.message],
      },
    }));
}
