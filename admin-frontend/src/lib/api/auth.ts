"use server";

import { Result } from "@/lib/types";

import { API_URL, resultify, sf } from "@/lib/utils";
import { cookies } from "next/headers";

import { JWT_KEY, setAuthToken } from "@/lib/jwt";
import { redirect } from "next/navigation";
import { z } from "zod";
import { jwtDecode } from "jwt-decode";
import { UserJwt } from "../types/auth";

const loginSchema = z.object({
  email: z.string().min(1, "Email is required"),
  password: z.string().min(1, "Password is required"),
});

export async function login(_initialState: any, formData: FormData): Promise<Result<void>> {
  // Parse form data
  const formObject = Object.fromEntries(formData);
  const { data, error } = loginSchema.safeParse(formObject);
  if (error) {
    return {
      error: error.issues[0].message,
    };
  }

  const res = await resultify(
    sf<string>(`${API_URL}/auth/login`, {
      method: "POST",
      body: JSON.stringify({
        email: data.email,
        password: data.password,
      }),
    })
  );

  // Handle errors
  if (res.error) {
    return res;
  }
  const jwt = jwtDecode(res.data!) as UserJwt;
  if (jwt.user.role != "admin") {
    return {
      error: "Unauthorized. You are not an admin",
    };
  }

  setAuthToken(res.data!);
  redirect("/");
}

export async function signout() {
  cookies().delete(JWT_KEY);
  redirect("/login");
}
