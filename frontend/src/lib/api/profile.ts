"use server";

import { UserProfileResponse } from "../interfaces/user";
import { Result } from "../types";
import { getAuthToken } from "../util/jwt";

export async function getProfile(): Promise<Result<UserProfileResponse>> {
  const jwt = getAuthToken();

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
