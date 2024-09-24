# Admin Dashboard

- **Endpoint:** `POST /disputes`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type SortOrder = "asc" | "desc";
type SortAttribute = "title" | "status" | "workflow" | "date_filed" | "date_resolved";

type FilterAttribute = "status" | "workflow";

interface Filter {
  // The attribute to filter by
  attr: FilterAttribute;

  // The value to search for.
  value: string;
}

interface AdminDisputesRequest {
  // Search term for the title of disputes
  search: string;

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
  filter: Filter[];

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
}>;

type AdminDisputesResponse = {
  disputes: AdminDisputes;
  total?: number;
};
```
