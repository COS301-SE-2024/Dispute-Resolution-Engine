export type SortOrder = "asc" | "desc";
export type SortAttribute = "title" | "status" | "workflow" | "date_filed" | "date_resolved";

export type FilterAttribute = "status" | "workflow";

export interface Filter {
  // The attribute to filter by
  attr: FilterAttribute;

  // The value to search for.
  value: string;
}

export interface AdminDisputesRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  sort?: {
    // The attribute to sort by
    attr: SortAttribute;

    // Sort order defaults to 'asc' if unspecified
    order?: SortOrder;
  };

  // The filters to apply to data
  filter?: Filter[];

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
  status: string;

  // The workflow that the dispute follows
  workflow: {
    id: string;
    title: string;
  };

  date_filed: string;

  // Optional because dispute may still be active (i.e. no resolved date)
  date_resolved?: string;
}

export type AdminDisputesResponse = Array<AdminDispute>;
