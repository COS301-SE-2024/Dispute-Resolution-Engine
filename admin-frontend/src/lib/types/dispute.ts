import { Filter, Sort } from ".";

export const DISPUTE_STATUS = [
  "Awaiting Respondant",
  "Active",
  "Review",
  "Appeal",
  "Settled",
  "Refused",
  "Withdrawn",
  "Transfer",
  "Other",
] as const;
export type DisputeStatus = (typeof DISPUTE_STATUS)[number];

export type DisputeFilter = Filter<"status" | "workflow">;
export type DisputeSort = Sort<"title" | "status" | "workflow" | "date_filed" | "date_resolved">;

export interface AdminDisputesRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  sort?: DisputeSort;

  // The filters to apply to data
  filter?: DisputeFilter[];

  dateFilter?: {
    filed?: {
      // Filter all disputes filed before the passed-in value (inclusive)
      before?: string;

      // Filter all disputes filed after the passed-in value (inclusive)
      after?: string;
    };

    // Specifying this filter would eliminate all unresolved disputes
    resolved?: {
      // Filter all disputes resolved before the passed-in value (inclusive)
      before?: string;

      // Filter all disputes resolved before the passed-in value (inclusive)
      after?: string;
    };
  };
}

export interface AdminDispute {
  id: string;
  title: string;
  status: DisputeStatus;

  // The workflow that the dispute follows
  workflow: {
    id: string;
    title: string;
  };

  date_filed: string;

  // Optional because dispute may still be active (i.e. no resolved date)
  date_resolved?: string;
}
export type AdminDisputesResponse = {
  disputes: Array<AdminDispute>;
  total: number;
};

export interface Evidence {
  id: string;
  url: string;
  label: string;
  submitted_at: string;
}

export interface UserDetails {
  name: string;
  email: string;
  address: string;
}

export interface DisputeDetails extends AdminDispute {
  description: string;
  evidence: Evidence[];
  complainant: UserDetails;
  respondent: UserDetails;
}

export type DisputeDetailsResponse = DisputeDetails;
