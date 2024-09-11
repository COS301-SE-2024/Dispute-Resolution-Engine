export const JWT_TIMEOUT = 60 * 20;
export const JWT_KEY = "jwt";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export function getAuthToken(): string {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    redirect("/login");
  }
  return jwt;
}

export function setAuthToken(jwt: string, timeout: number = JWT_TIMEOUT) {
  cookies().set(JWT_KEY, jwt, {
    maxAge: timeout,
    httpOnly: true,
  });
}
