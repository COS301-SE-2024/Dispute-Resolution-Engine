import {
  ArchiveGetResponse,
  ArchiveSearchRequest,
  ArchiveSearchResponse,
} from "../interfaces/archive";
import { Result } from "@/lib/types";

export async function searchArchive(
  params: ArchiveSearchRequest
): Promise<Result<ArchiveSearchResponse>> {
  // Uncomment once we have API
  //   return await fetch(`${API_URL}/disputes/search`, {
  //     method: "POST",
  //     body: JSON.stringify(params),
  //   }).then((res) => res.json());

  return {
    data: [...Array(params.limit ?? 10).keys()]
      .map((i) => i + (params.offset ?? 0))
      .map((i) => ({
        id: i.toString(),

        title: `Dispute #${i}`,
        summary:
          "Lorem ipsum dolor sit amet consectetur adipisicing elit. Soluta suscipit ducimus sequi alias tempora maxime odio libero delectus possimus aliquam ullam asperiores dolorem cumque, sunt numquam obcaecati. Eligendi, fugit commodi.",

        category: ["cooked", "cool"],

        date_filed: "yesterday",
        date_resolved: "today",

        resolution: "nah, i'd win",
      })),
  };
}
export async function fetchArchivedDispute(id: string): Promise<Result<ArchiveGetResponse>> {
  // Uncomment once we have API
  // return await fetch(`${API_URL}/disputes/archive/${id}`)
  //     .then((res) => res.json());

  return {
    data: {
      id: id,

      title: `Dispute #${id}`,
      summary: "Lorem Ipsum",

      category: ["cooked", "cool"],

      date_filed: "yesterday",
      date_resolved: "today",

      resolution: "nah, i'd win",
      events: [...Array(10).keys()].map((i) => ({
        timestamp: `${i}`,
        type: "uploaded",
        description: "documents uploaded",
      })),
    },
  };
}
