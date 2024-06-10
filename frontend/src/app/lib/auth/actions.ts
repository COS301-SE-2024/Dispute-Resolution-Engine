"use server";

import { Result } from "@/lib/types";
import { SignupError, signupSchema } from "./types";
import { redirect } from "next/navigation";
import { API_URL } from "@/lib/utils";

export async function signup(
  initialState: any,
  formData: FormData
): Promise<Result<string, SignupError>> {
  const { data, error } = signupSchema.safeParse({
    firstName: formData.get("firstName"),
    lastName: formData.get("lastName"),
    email: formData.get("email"),
    password: formData.get("password"),
    passwordConfirm: formData.get("passwordConfirm"),

    addrCountry: formData.get("addrCountry"),
    addrProvince: formData.get("addrProvince"),
    addrCity: formData.get("addrCity"),
    addrStreet3: formData.get("addrStreet3"),
    addrStreet2: formData.get("addrStreet2"),
    addrStreet: formData.get("addrStreet"),

    dateOfBirth: formData.get("dateOfBirth"),
    idNumber: formData.get("idNumber"),
  });

  if (error) {
    return {
      status: 500,
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
      password_hash: data.password,
      gender: "Male",
    }),
  }).then((res) => res.json());

  if (resError) {
    return {
      status: 500,
      error: {
        _errors: [resError],
      },
    };
  }
  return {
    status: 200,
    data: resData,
  };
}
