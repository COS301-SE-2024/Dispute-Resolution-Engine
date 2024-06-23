"use server";

import { z } from "zod";
import { DisputeCreateError, disputeCreateSchema } from "../schema/dispute";
import { Result } from "../types";

export async function createDispute(_initial: unknown, data: FormData): Promise<Result<string, DisputeCreateError>>  {
    const { data: parsed, error: parseErr } = disputeCreateSchema.safeParse(Object.fromEntries(data));
    if(parseErr) {
        return {
            error: parseErr.format()
        };
    }

    // Append each property of requestData to formData
    const formData = new FormData();
    formData.append("title", parsed.title);
    formData.append("description", parsed.summary);
    formData.append("respondent[full_name]", parsed.respondentName);
    formData.append("respondent[email]", parsed.respondentEmail);
    formData.append("respondent[telephone]", parsed.respondentTelephone);
    formData.append("files", formData.get("files")!);
    console.log(formData);

    return {
        data: ""
    }
}