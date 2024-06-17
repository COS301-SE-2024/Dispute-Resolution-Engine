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

// TODO: File upload endpoint
