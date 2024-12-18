import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { Result, resultSchema, TimerDuration } from "./types";

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
  return fetch(input, init).then(async (res) => {
    if (!res.ok) {
      let error = `Request failed with code ${res.status}`;
      try {
        const x = resultSchema.parse(await res.json());
        error = x.error ?? error;
      } catch (e) {}
      console.log("Response in sf", JSON.stringify(res));
      throw new Error(error);
    }
    return res;
  });
}

export function durationFromString(d: string): TimerDuration {
  const regex = /(\d+h)?(\d+m)?(\d+s)?/g;
  const dur = regex.exec(d);
  if (!dur) {
    throw new Error("bad boy");
  }
  let hours = dur[1] ? parseInt(dur[1].substring(0, dur[1].length - 1)) : 0;

  const days = Math.floor(hours / 24);
  hours = hours % 24;

  const minutes = dur[2] ? parseInt(dur[2].substring(0, dur[2].length - 1)) : 0;
  const seconds = dur[3] ? parseInt(dur[3].substring(0, dur[3].length - 1)) : 0;

  return { hours, days, minutes, seconds };
}

export function durationToString(d: TimerDuration): string {
  console.log("DUR", d);
  const hours = d.hours + d.days * 24;
  const { minutes, seconds } = d;

  let result = "";
  if (hours > 0) {
    result += `${hours}h`;
  }
  if (minutes > 0) {
    result += `${minutes}m`;
  }
  if (seconds > 0) {
    result += `${seconds}s`;
  }

  return result;
}
