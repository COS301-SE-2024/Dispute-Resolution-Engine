import { JWT_KEY } from "@/lib/constants";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  if (!cookies().get(JWT_KEY)) {
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: ["/disputes/:path*", "/profile", "/disputes"],
};
