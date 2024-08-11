# Archive

# Utility Types

```ts
interface ArchivedDisputeSummary {
  id: string;

  title: string;
  description: string;
  summary: string;

  category: string[];

  date_filed: string;
  date_resolved: string;

  resolution: string;
}
interface ArchivedDispute extends ArchivedDisputeSummary {
  events: {
    timestamp: string;
    type: string;
    description: string;
  }[];
}

type SortAttribute = "title" | "date_filed" | "date_resolved" | "date_filed" | "time_taken";
```

# Searching

- **Endpoint:** `POST /disputes/archive/search`

```ts
interface ArchiveSearchRequest {
  search: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  order?: "asc" | "desc";

  // What attribute to sort by
  sort?: SortAttribute;

  filter?: {
    category?: string[];
    time?: number;
  };
}
```

```ts
type ArchiveSearchResponse = ArchivedDisputeSummary[];
```

# Archived Dispute Details

- **Endpoint:** `GET /disputes/archive/{id}`

```ts
type ArchiveGetResponse = ArchivedDispute;
```
