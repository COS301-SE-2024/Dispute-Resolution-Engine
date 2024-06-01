"use server";

import { LoginData } from "../login/login-form";
import { SignupData } from "../signup/signup-form";

const API_URL = process.env.API_URL;

export async function signup(data: SignupData): Promise<string> {
  const response = await fetch(`${API_URL}/api`, {
    method: "POST",
    body: JSON.stringify({
      request_type: "create_account",
      body: {
        first_name: data.firstName,
        surname: data.lastName,
        email: data.email,
        password_hash: data.password
      }
    })
  }).then(res => res.json());
  if(typeof response == "string") {
    return response;
  }
  return response.Error;
}

export async function login(data: LoginData): Promise<string> {
  const response = await fetch(`${API_URL}/api`, {
    method: "POST",
    body: JSON.stringify({
      request_type: "login",
      body: {
        email: data.email,
        password_hash: data.password,
      }
    })
  }).then(res => res.json());
  if(typeof response == "string") {
    return response;
  }
  return response.Error;
}
