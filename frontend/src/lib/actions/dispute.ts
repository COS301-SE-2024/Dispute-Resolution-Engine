"use server";

import {
  DisputeCreateData,
  DisputeCreateError,
  ExpertRejectData,
  ExpertRejectError,
  disputeCreateSchema,
  expertRejectSchema,
} from "../schema/dispute";
import { Result } from "../types";
import { API_URL, formFetch } from "../utils";
import { cookies } from "next/headers";
import { JWT_KEY } from "../constants";

import { DisputeCreateResponse, DisputeEvidenceUploadResponse } from "../interfaces/dispute";
import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { getAuthToken } from "../util/jwt";

export async function createDispute(_initial: unknown, data: FormData): Promise<Result<string>> {
  const { data: parsed, error: parseErr } = disputeCreateSchema.safeParse(Object.fromEntries(data));
  if (parseErr) {
    return {
      error: parseErr.issues[0].message,
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

const res = await formFetch<DisputeCreateData, string>(${API_URL}/disputes/create, {
    method: "POST",
    headers: {
      // Sub this for the proper getAuthToken thing
      Authorization: Bearer ${getAuthToken()},
      },
      body: formData,
    },
  );
  if (res.error) {
    return {
      error: res.error._errors[0],
    };
  }

  revalidatePath("/disputes");
  revalidatePath("/disputes/create");
  revalidatePath(`/disputes/${res.data.id}`);
  redirect(`/disputes/${res.data.id}`);
}

export async function rejectExpert(
  _initial: unknown,
  data: FormData,
): Promise<Result<string, ExpertRejectError>> {
  const { data: parsed, error: parseErr } = expertRejectSchema.safeParse(Object.fromEntries(data));
  if (parseErr) {
    return {
      error: parseErr.format(),
    };
  }

  const res = await formFetch<ExpertRejectData, string>(
    `${API_URL}/disputes/${parsed.dispute_id}/experts/reject`,
    {
      method: "POST",
      headers: {
        // Sub this for the proper getAuthToken thing
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify({
        expert_id: parseInt(parsed.expert_id),
        reason: parsed.reason,
      }),
    },
  );

  if (!res.error) {
    revalidatePath(`/disputes/${parsed.dispute_id}`);
  }
  return res;
}

export async function uploadEvidence(
  _initial: unknown,
  data: FormData,
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
      Authorization: `Bearer ${getAuthToken()}`,
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
