/**
 * Generic Response from the API
 */

import { z } from "zod";

export const resultSchema = z.object({
  data: z.any().optional(),
  error: z.string().optional(),
});

export type Result<T, E = string> =
  | {
      data: T;
      error?: never;
    }
  | {
      data?: never;
      error: E;
    };

export interface Sort<T extends string> {
  // The attribute to sort by
  attr: T;

  // Sort order defaults to 'asc' if unspecified
  order?: "asc" | "desc";
}

export interface Filter<T extends string> {
  // The attribute to filter by
  attr: T;

  // The value to search for.
  value: string;
}

export * from "./dispute";
export * from "./workflow";
