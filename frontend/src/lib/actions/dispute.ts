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

  const res = await formFetch<DisputeCreateData, string>(`${API_URL}/disputes/create`, {
    method: "POST",
    headers: {
      // Sub this for the proper getAuthToken thing
      Authorization: `Bearer ${cookies().get(JWT_KEY)!.value}`,
    },
    body: formData,
  });
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

  const res = await formFetch<ExpertRejectData, string>(
    `${API_URL}/disputes/${parsed.dispute_id}/experts/reject`,
    {
      method: "POST",
      headers: {
        // Sub this for the proper getAuthToken thing
        Authorization: `Bearer ${cookies().get(JWT_KEY)!.value}`,
      },
      body: JSON.stringify({
        expert_id: parsed.expert_id,
        reason: parsed.reason,
      }),
    }
  );

  if (!res.error) {
    revalidatePath(`/disputes/${parsed.dispute_id}`);
  }
  return res;
}
