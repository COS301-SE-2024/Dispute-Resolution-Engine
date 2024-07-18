"use server";

import { DisputeCreateError, disputeCreateSchema } from "../schema/dispute";
import { Result } from "../types";
import { API_URL } from "../utils";
import { cookies } from "next/headers";
import { JWT_KEY } from "../constants";
import { DisputeEvidenceUploadResponse } from "../interfaces/dispute";
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
  data
    .getAll("file")
    .map((f) => f as File)
    .forEach((file) => formData.append("files", file, file.name));
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

export async function uploadEvidence(
  _initial: unknown,
  data: FormData
): Promise<Result<DisputeEvidenceUploadResponse>> {
  const disputeId = data.get("dispute_id");
  const formData = new FormData();
  data
    .getAll("files")
    .map((f) => f as File)
    .forEach((file) => formData.append("files", file, file.name));

  const res = await fetch(`${API_URL}/disputes/${disputeId}/evidence`, {
    method: "POST",
    headers: {
      // Sub this for the proper getAuthToken thing
      Authorization: `Bearer ${cookies().get(JWT_KEY)!.value}`,
    },
    body: formData,
  })
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));

  if (!res.error) {
    revalidatePath(`/dispute/${disputeId}`);
  }

  return res;
}
