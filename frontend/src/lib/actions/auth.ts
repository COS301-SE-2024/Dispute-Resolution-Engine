"use server";

import {
  LoginData,
  LoginError,
  ResetLinkData,
  ResetLinkError,
  ResetPassData,
  ResetPassError,
  loginSchema,
  resetLinkSchema,
  resetPassSchema,
  signupSchema,
  verifySchema,
} from "@/lib/schema/auth";
import { Result } from "@/lib/types";
import { API_URL, formFetch, resultify, safeFetch } from "@/lib/utils";
import { cookies } from "next/headers";

import { redirect } from "next/navigation";
import { JWT_KEY, JWT_VERIFY_TIMEOUT } from "../constants";
import { validateResult } from "../util";
import { getAuthToken, setAuthToken } from "../util/jwt";

export const signup = async (payload: unknown) => {
    const { data, error } = signupSchema.safeParse(payload);
    if (error) {
      throw new Error(error.issues[0].message);
    }

    // TODO: uncomment once API works
    const res = await fetch(`${API_URL}/auth/signup`, {
      method: "POST",
      body: JSON.stringify({
        first_name: data.firstName,
        surname: data.lastName,
        nationality: data.nationality,

        email: data.email,
        password: data.password,

        phone_number: "0110110110",
        birthdate: "2004-01-15",
        gender: "Male",
        timezone: ".",
        preferred_language: "English",
        user_type: "user",
      }),
    }).then(validateResult<string>);
    console.log("RECEVIED: ", res)
    setAuthToken(res, JWT_VERIFY_TIMEOUT);
    console.log("JWT SET: ")
    redirect("/signup/verify");
}

export async function login(
  _initialState: any,
  formData: FormData
): Promise<Result<string, LoginError>> {
  // Parse form data
  const formObject = Object.fromEntries(formData);
  const { data, error } = loginSchema.safeParse(formObject);
  if (error) {
    return {
      error: error.format(),
    };
  }

  // TODO: uncomment when API is working
  // Send request to the API
  const res = await formFetch<LoginData, string>(`${API_URL}/auth/login`, {
    method: "POST",
    body: JSON.stringify({
      email: data.email,
      password: data.password,
    }),
  });

  // Handle errors
  if (res.error) {
    return res;
  }
  setAuthToken(res.data);
  redirect("/disputes");
}
export async function signout() {
  cookies().delete(JWT_KEY);
  redirect("/login");
}

export const verify = (payload: unknown) =>
  resultify(async () => {
    // Parse the form data
    const { data, error } = verifySchema.safeParse(payload);
    if (error) {
      throw new Error(error.issues[0].message);
    }

    // TODO: uncomment once API works
    // Send request to the API
    const res = await fetch(`${API_URL}/auth/verify`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify({
        pin: data.pin,
      }),
    }).then(validateResult<string>);

    // Everything good
    setAuthToken(res);
    redirect("/disputes");
  });

export async function resendOTP(_initialState: any, _: FormData): Promise<Result<string>> {
  // Retrieve the temporary JWT

  // TODO: uncomment once API works
  const res = await safeFetch<string>(`${API_URL}/auth/resend-otp`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${getAuthToken()}`,
    },
  });

  return res;
}

export async function sendResetLink(
  _initialState: any,
  formData: FormData
): Promise<Result<string, ResetLinkError>> {
  // Parse form data
  const formObject = Object.fromEntries(formData);
  const { data, error } = resetLinkSchema.safeParse(formObject);
  if (error) {
    return {
      error: error.format(),
    };
  }

  const res = await formFetch<ResetLinkData, string>(`${API_URL}/auth/reset-password/send-email`, {
    method: "POST",
    body: JSON.stringify({
      email: data.email,
    }),
  });

  // Handle errors
  if (res.error) {
    return res;
  }
  setAuthToken(res.data);
  redirect("/reset/success");
}

export async function resetPassword(
  _initialState: any,
  formData: FormData
): Promise<Result<string, ResetPassError>> {
  // Parse form data
  const formObject = Object.fromEntries(formData);
  const { data, error } = resetPassSchema.safeParse(formObject);
  if (error) {
    return {
      error: error.format(),
    };
  }

  const res = await formFetch<ResetPassData, string>(`${API_URL}/auth/reset-password/reset`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${data.jwt}`,
    },
    body: JSON.stringify({
      newPassword: data.password,
    }),
  });

  // Handle errors
  if (res.error) {
    return res;
  }
  setAuthToken(res.data);
  redirect("/login");
}
