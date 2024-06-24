import {
  ArchiveGetResponse,
  ArchiveSearchRequest,
  ArchiveSearchResponse,
} from "../interfaces/archive";
import { Result } from "@/lib/types";
import { API_URL } from "../utils";

export async function searchArchive(
  params: ArchiveSearchRequest
): Promise<Result<ArchiveSearchResponse>> {
  const res = await fetch(`${API_URL}/disputes/archive/search`, {
    method: "POST",
    body: JSON.stringify(params),
  });
  return res.json().catch(async (e: Error) => ({
    error: e.message,
  }));
}
export async function fetchArchivedDispute(id: string): Promise<Result<ArchiveGetResponse>> {
  return fetch(`${API_URL}/disputes/archive/${id}`)
    .then((res) => res.json())
    .catch((e: Error) => ({
      error: e.message,
    }));
}
