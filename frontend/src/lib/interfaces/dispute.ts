export type Role = "Complainant" | "Respondent";

export const DISPUTE_STATUS = [
  "Awaiting Respondant",
  "Active",
  "Review",
  "Settled",
  "Refused",
  "Withdrawn",
  "Transfer",
  "Appeal",
  "Other",
] as const;
export type DisputeStatus = (typeof DISPUTE_STATUS)[number];

export const DISPUTE_DECISION = ["Settled", "Refused", "Appeal"] as const;
export type DisputeDecision = (typeof DISPUTE_DECISION)[number];

export type DisputeSummary = {
  id: string;
  title: string;
  description: string;
  status: DisputeStatus;
  role: Role;
};

export type Evidence = {
  label: string;
  url: string;
  date_submitted: string;
  uploader_role: Role;
};

export type Expert = {
  id: string;
  role: string;
  full_name: string;
  email: string;
  phone: string;
};

export type DisputeListResponse = DisputeSummary[];
export type DisputeResponse = {
  id: string;
  title: string;
  description: string;
  status: string;
  case_date: string;
  role: Role;

  evidence: Evidence[];
  experts: Expert[];
};

/**
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 *
 * > **Jurisdiction Agreement**:
 * > - Agreement to submit to the jurisdiction of the High Court of South Africa.
 *
 * What does this mean???
 *
 * This needs to be refactored to include FormData
 */
export interface DisputeCreateRequest {
  title: string;
  description: string;
  evidence: File[];
  respondent: {
    full_name: string;
    email: string;
    telephone: string;
  };
}

export interface DisputeCreateResponse {
  id: number;
}

export interface DisputeStatusUpdateRequest {
  dispute_id: string;
  status: string;
}
// TODO: File upload endpoint
export type DisputeEvidenceUploadResponse = string;
