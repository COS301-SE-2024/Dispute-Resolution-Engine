"use server";

import { Result } from "@/lib/types";
import {
  LoginError,
  SignupError,
  loginSchema,
  signupSchema,
  verifySchema,
} from "@/lib/schema/auth";
import { API_URL } from "@/lib/utils";
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
  const resData = "unverified";
  const resError = undefined;

  // TODO: uncomment once API works
  // const { data: resData, error: resError } = await fetch(`${API_URL}/auth/signup`, {
  //   method: "POST",
  //   body: JSON.stringify({
  //     first_name: data.firstName,
  //     surname: data.lastName,
  //     email: data.email,
  //     phone_number: "0110110110",

  //     password: data.password,

  //     birthdate: data.dateOfBirth,
  //     gender: "Male",
  //     nationality: "za",

  //     timezone: ".",
  //     preferred_language: ".",
  //   }),
  // }).then((res) => res.json());

  if (resError) {
    return {
      error: {
        _errors: [resError],
      },
    };
  }
  setAuth(resData, JWT_VERIFY_TIMEOUT);
  redirect("/signup/verify");
  return {
    data: "Signup successful",
  };
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

  const resData = "working";
  const resError = undefined;

  // TODO: uncomment when API is working
  // Send request to the API
  // const { data: resData, error: resError } = await fetch(`${API_URL}/auth/login`, {
  //   method: "POST",
  //   body: JSON.stringify({
  //     email: data.email,
  //     password: data.password,
  //   }),
  // }).then((res) => res.json());
  // console.log(resData);

  // Handle errors
  if (resError) {
    return {
      error: {
        _errors: [resError],
      },
    };
  }
  setAuth(resData);
  return {
    data: "Login successful",
  };
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
        _errors: ["OTP Expired"],
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

  const resData = "working";
  const resError = data.pin == "123456" ? undefined : "Incorrect OTP";

  // TODO: uncomment once API works
  // Send request to the API
  // const { data: resData, error: resError } = await fetch(`${API_URL}/auth/verify`, {
  //   method: "POST",
  //   headers: {
  //     Authorization: `Bearer ${jwt}`,
  //   },
  //   body: JSON.stringify({
  //     email: data.email,
  //     password: data.password,
  //   }),
  // }).then((res) => res.json());
  // console.log(resData);

  // Handle Errors
  if (resError) {
    return {
      error: {
        _errors: [resError],
      },
    };
  }

  // Everything good
  setAuth(resData);
  return {
    data: "Email verified and logged in",
  };
}
