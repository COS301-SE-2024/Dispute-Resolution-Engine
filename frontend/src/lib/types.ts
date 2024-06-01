/**
 * Generic Response from the API
 */
export type Result<T, E = string> =
  | {
      status: number;
      data: T;
      error?: never;
    }
  | {
      status: number;
      data?: never;
      error: E;
    };
