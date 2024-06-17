export type DisputeSummary = {
  id: string;
  title: string;
  description: string;
  status: string;
};

export type Evidence = {
  label: string;
  url: string;
  date_submitted: string;
};

export type Expert = {
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
  date_created: string;

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
  description: string;
  desired_outcome: string;

  respondent: {
    full_name: string;
    email: string;
    telephone: string;
  };

  jurisdictional_basis: Evidence;

  /**
   * IDs of all adjudicators to be appointed
   */
  adjudicators: string[];

  // This should be FormData, but I don't know how to annotate that
  evidence: Evidence[];
}

export type DisputeCreateResponse = string;

// TODO: File upload endpoint
