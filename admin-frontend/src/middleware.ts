import { JWT_KEY } from "@/lib/jwt";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  if (request.nextUrl.pathname !=="/login" && !cookies().get(JWT_KEY)) {
    return NextResponse.redirect(new URL("/admin/login", request.url));
  }
}

// See "Matching Paths" below to learn more
// export const config = {
//   matcher: ["/:path*"],
// };
