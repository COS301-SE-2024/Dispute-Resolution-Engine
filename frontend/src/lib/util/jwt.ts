"use server";

import { cookies } from "next/headers";
import { JWT_KEY } from "../constants";
import { redirect } from "next/navigation";

export function getAuthToken(): string {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    redirect("/login");
  }
  return jwt;
}
