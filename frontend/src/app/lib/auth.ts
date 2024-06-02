"use server";

import { Result } from "@/lib/types";
import { LoginData } from "../login/login-form";
import { SignupData } from "../signup/signup-form";

const API_URL = process.env.API_URL;

export async function signup(data: SignupData): Promise<Result<{ message: string }>> {
  const response = await fetch(`${API_URL}/api`, {
    cache: "no-cache",
    method: "POST",
    body: JSON.stringify({
      request_type: "create_account",
      body: {
        first_name: data.firstName,
        surname: data.lastName,
        email: data.email,
        password: data.password,
      },
    }),
  }).then((res) => res.json());
  return response;
}

export async function login(data: LoginData): Promise<Result<{ message: string }>> {
  const response = await fetch(`${API_URL}/api`, {
    cache: "no-cache",
    method: "POST",
    body: JSON.stringify({
      request_type: "login",
      body: {
        email: data.email,
        password: data.password,
      },
    }),
  }).then((res) => res.json());
  console.log(response);

  return response;
}
