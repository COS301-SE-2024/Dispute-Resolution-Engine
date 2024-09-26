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
