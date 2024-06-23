"use server";

import { UserProfileResponse } from "../interfaces/user";
import { Result } from "../types";
import { getAuthToken } from "../util/jwt";
import { API_URL } from "../utils";

export async function getProfile(): Promise<Result<UserProfileResponse>> {
  const jwt = getAuthToken();

  // TODO: Uncomment this once API works
  return fetch(`${API_URL}/user/profile`, {
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
  })
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
}
