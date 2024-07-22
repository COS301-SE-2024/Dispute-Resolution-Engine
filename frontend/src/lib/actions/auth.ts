"use server";

import { Result } from "@/lib/types";
import {
  LoginData,
  LoginError,
  ResetLinkData,
  ResetLinkError,
  ResetPassData,
  ResetPassError,
  SignupData,
  SignupError,
  loginSchema,
  resetLinkSchema,
  resetPassSchema,
  signupSchema,
  verifySchema,
} from "@/lib/schema/auth";
import { API_URL, formFetch, safeFetch } from "@/lib/utils";
import { cookies } from "next/headers";

import { JWT_TIMEOUT, JWT_KEY, JWT_VERIFY_TIMEOUT } from "../constants";
import { redirect } from "next/navigation";

function setAuth(jwt: string, timeout: number = JWT_TIMEOUT) {
  cookies().set(JWT_KEY, jwt, {
    maxAge: timeout,
    httpOnly: true,
  });
}
function getAuth() {
  return cookies().get(JWT_KEY)?.value;
}

export async function signup(
  _initialState: any,
  formData: FormData
): Promise<Result<string, SignupError>> {
  const formObject = Object.fromEntries(formData);
  const { data, error } = signupSchema.safeParse(formObject);

  if (error) {
    return {
      error: error.format(),
    };
  }

  // TODO: uncomment once API works
  const res = await formFetch<SignupData, string>(`${API_URL}/auth/signup`, {
    method: "POST",
    body: JSON.stringify({
      first_name: data.firstName,
      surname: data.lastName,
      email: data.email,
      phone_number: "0110110110",

      password: data.password,

      birthdate: data.dateOfBirth,
      gender: data.gender,
      nationality: data.nationality,

      timezone: ".",
      preferred_language: data.preferredLanguage,
    }),
  });

  if (res.error) {
    return res;
  }

  setAuth(res.data, JWT_VERIFY_TIMEOUT);
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
  setAuth(res.data);
  redirect("/disputes");
}

export async function verify(
  _initialState: any,
  formData: FormData
): Promise<Result<string, LoginError>> {
  // Retrieve the temporary JWT
  const jwt = getAuth();
  if (!jwt) {
    return {
      error: {
        _errors: ["JWT Expired"],
      },
    };
  }

  // Parse the form data
  const formObject = Object.fromEntries(formData);
  const { data, error } = verifySchema.safeParse(formObject);
  if (error) {
    return {
      error: error.format(),
    };
  }

  // TODO: uncomment once API works
  // Send request to the API
  const res = await formFetch<LoginData, string>(`${API_URL}/auth/verify`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
    body: JSON.stringify({
      pin: data.pin,
    }),
  });

  // Handle Errors
  if (res.error) {
    return res;
  }

  // Everything good
  setAuth(res.data);
  redirect("/disputes");
  return {
    data: "Email verified and logged in",
  };
}

export async function resendOTP(_initialState: any, formData: FormData): Promise<Result<string>> {
  // Retrieve the temporary JWT
  const jwt = getAuth();
  if (!jwt) {
    return {
      error: "JWT Expired",
    };
  }

  // TODO: uncomment once API works
  const res = await safeFetch<string>(`${API_URL}/auth/resend-otp`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${jwt}`,
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
  setAuth(res.data);
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
  setAuth(res.data);
  redirect("/login");
}
