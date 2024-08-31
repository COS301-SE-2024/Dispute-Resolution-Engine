import { Result } from "@/lib/types";
import { AdminDisputesRequest, AdminDisputesResponse } from "../types/dispute";

export async function getDisputeList(
  req: AdminDisputesRequest
): Promise<Result<AdminDisputesResponse>> {
  return {
    data: [
      {
        id: "1",
        title: "Sales hire vs. Hire City",
        status: "Awaiting Respondent",

        workflow: {
          id: "1",
          title: "Domain Dispute",
        },

        date_filed: "2 days ago",
      },
      {
        id: "2",
        title: "Telkom SA vs. Cool-aid",
        status: "Active",

        // The workflow that the dispute follows
        workflow: {
          id: "2",
          title: "Marital Dispute",
        },

        date_filed: "yesterday",
      },
      {
        id: "3",
        title: "Standard Bank vs. Bank of Standards",
        status: "Settled",

        // The workflow that the dispute follows
        workflow: {
          id: "1",
          title: "Domain Dispute",
        },

        date_filed: "2 days ago",
        date_resolved: "today",
      },
    ],
  };
}
