"use server";

import { LoginData } from "../login/login-form";
import { SignupData } from "../signup/signup-form";

export async function signup(data: SignupData): Promise<string> {
  return new Promise((resolve) =>
    setTimeout(() => {
      console.log(data);
      resolve("Something happened");
    }, 2000)
  );
}

export async function login(data: LoginData): Promise<string> {
  return new Promise((resolve) =>
    setTimeout(() => {
      console.log(data);
      resolve("Logged in!!");
    }, 2000)
  );
}
