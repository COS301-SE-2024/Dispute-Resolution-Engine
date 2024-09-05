"use server";

import { Result } from "@/lib/types";
import {
  AdminDisputesRequest,
  AdminDisputesResponse,
  DisputeDetails,
  DisputeDetailsResponse,
  DisputeStatus,
} from "@/lib/types";

export async function getDisputeList(
  req: AdminDisputesRequest
): Promise<Result<AdminDisputesResponse>> {
  return {
    data: MOCK_DATA,
  };
}

export async function getDisputeDetails(id: string): Promise<Result<DisputeDetailsResponse>> {
  const result = MOCK_DATA.find((d) => d.id === id);
  if (!result) {
    return {
      error: "Dispute not found",
    };
  } else {
    return {
      data: result,
    };
  }
}

export async function changeDisputeStatus(
  id: string,
  status: DisputeStatus
): Promise<Result<string>> {
  return {
    data: "Status changed",
    // error: "Big bad error things",
  };
}
export async function deleteEvidence(
  disputeId: string,
  evidenceId: string
): Promise<Result<string>> {
  return {
    data: "Evidence deleted",
  };
}

const MOCK_DATA: DisputeDetails[] = [
  {
    id: "ZA2007-0001",
    title: "Sales hire vs. Hire City",
    status: "Awaiting Respondent",

    workflow: {
      id: "1",
      title: "Domain Dispute",
    },

    date_filed: "2 days ago",
    description:
      "The Complainant contends that it has rights in respect of the name or  mark MR PLASTIC and that the domain name in dispute is identical or  similar to this name or mark and it is therefore an abusive registration.",
    evidence: [],
    complainant: {
      name: "Mr. Plastic CC",
      email: "mrplastic@gmail.com",
      address: "13 Geldenhuys Road\nMalvern East\nBedfordview, Gauteng",
    },
    respondent: {
      name: "Mr.  Plastic  &  Mining  Promotional Goods",
      email: "mrplastic@outlook.com",
      address: "26 Boom Street\nJeppestown, Gauteng",
    },
  },
  {
    id: "ZA2007-0003",
    title: "Telkom SA LTD vs. Cool Ideas 1290 CC",
    status: "Awaiting Respondent",

    workflow: {
      id: "1",
      title: "Domain Dispute",
    },

    date_filed: "2 days ago",
    description:
      "It has registered trade mark rights. It has listed 10 (ten) trade mark  registrations in South Africa dating from 1991 for the trade mark TELKOM and  TELKOM & KEYPAD logo in various classes including class 38 that relates to  telecommunication services.",

    evidence: [],
    respondent: {
      name: "Cool Ideas 1290 CC",
      email: "telkom@gmail.com",
      address: "25 Sidonia Avenue\nNorwood, Johannesburg\nGAUTENG",
    },
    complainant: {
      name: "Telkom SA Limited",
      email: "telkom@outlook.com",
      address: "Telkom Towers North\n152 Proes Street\nPretoria\nGauteng",
    },
  },
];
