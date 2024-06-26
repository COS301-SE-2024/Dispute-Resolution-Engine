import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export type Slottable<P> = P & {
  asChild?: boolean;
};

export const API_URL = process.env.API_URL;
