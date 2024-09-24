import { JWT_KEY } from "@/lib/jwt";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { jwtDecode } from "jwt-decode";
import { UserJwt } from "./lib/types/auth";

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  // Forces users to login
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // TODO: add some measure to prevent a malicious actor from setting the
  // user role to "admin", bypassing this rule. One potential solution is
  // a UDP lookup to the API to check if the JWT is valid or not.
  const data = jwtDecode(jwt) as UserJwt;
  if (data.user.role != "admin") {
    cookies().delete(JWT_KEY);
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: ["/disputes(.*)", "/experts(.*)", "/workflows(.*)", "/tickets(.*)"],
};
