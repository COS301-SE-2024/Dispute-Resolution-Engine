"use server";

import { Result } from "@/lib/types";
import { SignupError, signupSchema } from "./types";
import { redirect } from "next/navigation";

export async function signup2(
  initialState: any,
  formData: FormData
): Promise<Result<string, SignupError>> {
  const data = signupSchema.safeParse({
    firstName: formData.get("firstName"),
    lastName: formData.get("lastName"),
    email: formData.get("email"),
    password: formData.get("password"),
    passwordConfirm: formData.get("passwordConfirm"),
  });
  if (data.error) {
    return {
      status: 500,
      error: data.error.format(),
    };
  }

  return new Promise((res) => {
    setTimeout(() => {
      res({
        status: 200,
        data: "Signup successful",
      });
      redirect("/disputes");
    }, 2000);
  });
}
