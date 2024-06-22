"use server";

import { cookies } from "next/headers";
import { UserProfileResponse } from "../interfaces/user";
import { Result } from "../types";
import { JWT_KEY } from "../constants";
import { API_URL } from "../utils";

export async function getProfile(): Promise<Result<UserProfileResponse>> {
  const jwt = cookies().get(JWT_KEY)?.value;
  if (!jwt) {
    return {
      error: "Unauthorized",
    };
  }

  return {
    data: {
      first_name: "John",
      surname: "Doe",
      email: "user@example.com",
      phone_number: "0-11",

      birthdate: "today",
      gender: "Female",
      nationality: "South africa",

      timezone: "today",
      preferred_language: "en-US",

      addresses: [],
      useDarkMode: true,
    },
  };

  // TODO: Uncomment this once API works
  //   return fetch(`${API_URL}/user/profile`, {
  //     headers: {
  //       Authorization: `Bearer ${jwt}`,
  //     },
  //   })
  //     .then((res) => res.json())
  //     .catch((e: Error) => ({
  //       error: e.message,
  //     }));
}
