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

export * from "./dispute";
