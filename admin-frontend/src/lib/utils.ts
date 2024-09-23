import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { Result } from "./types";

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

export const API_URL = process.env.API_URL;
