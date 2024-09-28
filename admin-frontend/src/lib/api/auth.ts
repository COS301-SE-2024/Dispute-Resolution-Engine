"use server";

import { Result } from "@/lib/types";

import { API_URL, resultify, sf, validateResult } from "@/lib/utils";
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

export async function login(_initialState: any, formData: FormData): Promise<Result<never>> {
  return resultify(
    (async () => {
      // Parse form data
      const formObject = Object.fromEntries(formData);
      const { data, error } = loginSchema.safeParse(formObject);
      if (error) {
        throw new Error(error.issues[0].message);
      }

      const res = await sf(`${API_URL}/auth/login`, {
        method: "POST",
        body: JSON.stringify({
          email: data.email,
          password: data.password,
        }),
      }).then(validateResult<string>);

      const jwt = jwtDecode(res) as UserJwt;
      if (jwt.user.role != "admin") {
        throw new Error("Unauthorized. You are not an admin");
      }

      setAuthToken(res);
      redirect("/");
    })()
  );
}

export async function signout() {
  cookies().delete(JWT_KEY);
  redirect("/login");
}
