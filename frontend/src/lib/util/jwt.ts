import { cookies } from "next/headers";
import { JWT_KEY, JWT_TIMEOUT } from "../constants";
import { redirect } from "next/navigation";

export function getAuthToken(): string | never {
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
