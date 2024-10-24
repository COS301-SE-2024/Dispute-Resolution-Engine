import { JWT_KEY } from "@/lib/jwt";
import { jwtDecode } from "jwt-decode";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { UserJwt } from "./lib/types/auth";

// Helper function to check if the path matches any in the validPaths array
function isPathValid(path: string, validPaths: string[]) {
  return validPaths.some((validPath) => {
    const pathRegex = new RegExp(`^${validPath.replace("*", ".*")}$`);
    return pathRegex.test(path);
  });
}

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  if (request.nextUrl.pathname == "/admin/login" || request.nextUrl.pathname.startsWith("/_next")) {
    return;
  }

  const cookie = cookies().get(JWT_KEY);
  if (!cookie) {
    return NextResponse.rewrite(new URL("/admin/login", request.url));
  }

  const user = jwtDecode(cookie.value) as UserJwt;
  if (user.user.role != "admin") {
    return NextResponse.rewrite(new URL("/admin/login", request.url));
  }

  // if (cookies().get(JWT_KEY)) console.log(request.nextUrl.pathname);
  // const validPaths = ["/tickets*", "/workflows*", "/disputs*", "/experts*"];

  // const { pathname } = request.nextUrl;

  // // Check if the path matches any validPaths and the JWT cookie is not present
  // if (isPathValid(pathname, validPaths) && !cookies().get(JWT_KEY)) {
  //   // Rewrite the request to the login page
  //   return NextResponse.rewrite(new URL("/admin/login", request.url));
  // }

  // // Continue with the request if the conditions are not met
  // return NextResponse.next();
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: ["/:path*"],
};
