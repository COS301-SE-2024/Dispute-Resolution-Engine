"use server";

import { SignupData } from "../signup/signup-form";

export async function signup(data: SignupData): Promise<string> {
  return new Promise((resolve) =>
    setTimeout(() => {
      console.log(data);
      resolve("Something happened");
    }, 2000)
  );
}
