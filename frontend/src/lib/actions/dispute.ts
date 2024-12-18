"use server";

import {
  DisputeCreateData,
  DisputeDecisionError,
  ExpertRejectData,
  ExpertRejectError,
  disputeCreateSchema,
  disputeDecisionSchema,
  expertRejectSchema,
} from "../schema/dispute";
import { Result } from "../types";
import { API_URL, formFetch, resultify } from "../utils";

import { DisputeCreateResponse, DisputeEvidenceUploadResponse } from "../interfaces/dispute";
import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { getAuthToken } from "../util/jwt";
import { sf, validateResult } from "../util";

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
  formData.append("respondent[workflow]", parsed.workflow);
  data
    .getAll("file")
    .map((f) => f as File)
    .forEach((file) => formData.append("files", file, file.name));
  console.log(formData);

  const res = await formFetch<DisputeCreateData, DisputeCreateResponse>(
    `${API_URL}/disputes/create`,
    {
      method: "POST",
      headers: {
        // Sub this for the proper getAuthToken thing
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: formData,
    }
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

export const rejectExpert = (data: unknown) =>
  resultify(async () => {
    const { data: parsed, error: parseErr } = expertRejectSchema.safeParse(data);
    if (parseErr) {
      throw new Error(parseErr.issues[0].message);
    }

    const res = await sf(`${API_URL}/disputes/${parsed.dispute_id}/objections`, {
      method: "POST",
      headers: {
        // Sub this for the proper getAuthToken thing
        Authorization: `Bearer ${getAuthToken()}`,
      },
      body: JSON.stringify({
        expert_id: parseInt(parsed.expert_id),
        reason: parsed.reason,
      }),
    }).then(validateResult<number>);

    revalidatePath(`/disputes/${parsed.dispute_id}`);
    return res;
  });

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

function fileExists(data: FormData, key: string): boolean {
  if (!(data.get(key) instanceof File)) {
    return false;
  }

  const file = data.get("writeup") as File;
  return !(file.size === 0 || file.name === "undefined");
}

export async function uploadDecision(
  _initial: unknown,
  data: FormData
): Promise<Result<string, DisputeDecisionError>> {
  let { data: parsed, error: parseErr } = disputeDecisionSchema.safeParse(Object.fromEntries(data));

  if (!fileExists(data, "writeup")) {
    let error = parseErr?.format();
    return {
      error: {
        ...error,
        _errors: error?._errors ?? [],
        writeup: {
          _errors: ["Missing Writeup"],
        },
      },
    };
  }

  if (parseErr) {
    return { error: parseErr.format() };
  }

  const res = await fetch(`${API_URL}/disputes/${parsed?.dispute_id}/decision`, {
    method: "POST",
    headers: {
      // Sub this for the proper getAuthToken thing
      Authorization: `Bearer ${getAuthToken()}`,
    },
    body: data,
  })
    .then((res) => {
      if (!res.ok) {
        throw new Error(`Request failed with code ${res.status}`);
      }
      return res;
    })
    .then(() => ({ data: null } as Result<null>))
    .catch((e) => ({ error: e.message } as Result<null>));

  if (!res.error) {
    revalidatePath(`/dispute/${parsed?.dispute_id}`);
  }

  return {
    data: "",
  };
}
