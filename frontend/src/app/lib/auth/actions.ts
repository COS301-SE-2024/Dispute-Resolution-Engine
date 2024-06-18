"use server";

import { Result } from "@/lib/types";
import { LoginError, SignupError, loginSchema, signupSchema } from "./types";
import { API_URL } from "@/lib/utils";

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

  const { data: resData, error: resError } = await fetch(`${API_URL}/createAcc`, {
    method: "POST",
    body: JSON.stringify({
      first_name: data.firstName,
      surname: data.lastName,
      birthdate: data.dateOfBirth,
      nationality: data.addrCountry,
      email: data.email,
      password: data.password,
      gender: "Male",
    }),
  }).then((res) => res.json());

  if (resError) {
    return {
      error: {
        _errors: [resError],
      },
    };
  }
  return {
    data: resData,
  };
}

export async function login(
  _initialState: any,
  formData: FormData
): Promise<Result<string, LoginError>> {
  const formObject = Object.fromEntries(formData);
  const { data, error } = loginSchema.safeParse(formObject);

  if (error) {
    return {
      error: error.format(),
    };
  }

  const { data: resData, error: resError } = await fetch(`${API_URL}/login`, {
    method: "POST",
    body: JSON.stringify({
      email: data.email,
      password: data.password,
    }),
  }).then((res) => res.json());
  console.log(resData);

  if (resError) {
    return {
      error: {
        _errors: [resError],
      },
    };
  }
  return {
    data: resData,
  };
}
