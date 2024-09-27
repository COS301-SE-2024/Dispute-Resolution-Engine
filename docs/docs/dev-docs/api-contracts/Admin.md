# Admin Dashboard

## Helper Types

```ts
export interface Sort<T extends string> {
  // The attribute to sort by
  attr: T;

  // Sort order defaults to 'asc' if unspecified
  order?: "asc" | "desc";
}

export interface Filter<T extends string> {
  // The attribute to filter by
  attr: T;

  // The value to search for.
  value: string;
}

type ExpertStatus = "Approved"|"Rejected"|"Review"
export interface ExpertSummary {
  id: number;
  fullname: string;
  status: ExpertStatus;
}

```

# Disputes

- **Endpoint:** `POST /disputes`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type SortAttribute = "title" | "status" | "workflow" | "date_filed" | "date_resolved";
type FilterAttribute = "status" | "workflow";

interface AdminDisputesRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  sort?: Sort<SortAttribute>;

  // The filters to apply to data
  filter?: Filter<FilterAttribute>[];

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
```

The response will be an array of disputes:

```ts
type AdminDisputes = Array<{
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
  experts: []ExpertSummary;
}>;

type AdminDisputesResponse = {
  data: AdminDisputes;
  total: number;
};
```
