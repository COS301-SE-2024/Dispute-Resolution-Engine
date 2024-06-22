"use server";

import { UserProfileUpdateResponse } from "../interfaces/user";
import { ProfileError, profileSchema } from "../schema/profile";
import { Result } from "../types";
import { getAuthToken } from "../util/jwt";
import { API_URL } from "../utils";

export async function updateProfile(
  _init: unknown,
  formData: FormData
): Promise<Result<UserProfileUpdateResponse, ProfileError>> {
  const { data: parsed, error: parseErr } = profileSchema.safeParse(Object.fromEntries(formData));
  if (parseErr) {
    return {
      error: parseErr.format(),
    };
  }

  const { data, error } = await fetch(`${API_URL}/profile`, {
    method: "PATCH",
    headers: {
      Authorization: getAuthToken(),
    },
    body: JSON.stringify(parsed),
  })
    .then((r) => r.json())
    .catch((e: Error) => ({
      error: e.message,
    }));

  if (error) {
    return {
      error: {
        _errors: [error],
      },
    };
  }
  return {
    data,
  };
}
