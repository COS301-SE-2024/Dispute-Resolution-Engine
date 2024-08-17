/**
 * Generic Response from the API
 */
export type Result<T, E = string> =
  | {
      data: T;
      error?: never;
    }
  | {
      data?: never;
      error: E;
    };
