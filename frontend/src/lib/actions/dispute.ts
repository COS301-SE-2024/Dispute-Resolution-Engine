"use server";

import { DisputeCreateError, ExpertRejectError, disputeCreateSchema, expertRejectSchema } from "../schema/dispute";
import { Result } from "../types";
import { API_URL } from "../utils";
import { cookies } from "next/headers";
import { JWT_KEY } from "../constants";
import { revalidatePath } from "next/cache";

export async function createDispute(
  _initial: unknown,
  data: FormData
): Promise<Result<string, DisputeCreateError>> {
  const { data: parsed, error: parseErr } = disputeCreateSchema.safeParse(Object.fromEntries(data));
  if (parseErr) {
    return {
      error: parseErr.format(),
    };
  }

  // Append each property of requestData to formData
  const formData = new FormData();
  formData.append("title", parsed.title);
  formData.append("description", parsed.summary);
  formData.append("respondent[full_name]", parsed.respondentName);
  formData.append("respondent[email]", parsed.respondentEmail);
  formData.append("respondent[telephone]", parsed.respondentTelephone);
  formData.append("files", data.get("file")!);
  console.log(formData);

  const res = await fetch(`${API_URL}/disputes/create`, {
    method: "POST",
    headers: {
      // Sub this for the proper getAuthToken thing
      Authorization: `Bearer ${cookies().get(JWT_KEY)!.value}`,
    },
    body: formData,
  })
    .then((res) => res.json())
    .then((res) =>
      !res.error
        ? res
        : {
            error: {
              _errors: [res.error],
            },
          }
    )
    .catch((e: Error) => ({
      error: {
        _errors: [e.message],
      },
    }));
  console.log(res);
  return res;
}

export async function rejectExpert(
  _initial: unknown,
  data: FormData
): Promise<Result<string, ExpertRejectError>> {
  const { data: parsed, error: parseErr } = expertRejectSchema.safeParse(Object.fromEntries(data));
  if (parseErr) {
    return {
      error: parseErr.format(),
    };
  }

  revalidatePath(`/disputes/${parsed.dispute_id}`);

  return {
      data: "wauw"
  }
}
