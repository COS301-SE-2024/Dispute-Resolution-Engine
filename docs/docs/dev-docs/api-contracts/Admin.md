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
}>;

type AdminDisputesResponse = {
  data: AdminDisputes;
  total?: int;
};
```

# Tickets

The follwoing types are used to describe tickets:

```ts
interface TicketSummary {
  id: string;

  // The user that created the ticket
  user: {
    id: string;
    full_name: string;
  };

  // Timestamp of when the ticket was created
  date_created: string;

  // Timestamp of when the ticket was created
  subject: string;

  // The status of the ticket
  status: TicketStatus;
}

interface Ticket extends TicketSummary {
  // The initial message submitted with the ticket
  body: string;

  // All messages exchanged in the ticket (sorted by date)
  messages: TicketMessage[];
}

interface TicketMessage {
  id: string;

  // The user that submitted the message
  user: {
    id: string;
    full_name: string;
  };

  // The timestamp when the user submitted the ticket
  date_sent: string;

  // The message in the ticket
  message: string;
}
```

## Ticket List

- **Endpoint:** `POST /tickets`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type SortAttribute =
  | "date_created" // The date is was created
  | "subject" // The subject of the ticket
  | "status" // The status of the ticket
  | "user"; // The full name of the user;

type FilterAttribute = "status";

interface TicketListRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  sort?: Sort<SortAttribute>;

  // The filters to apply to data
  filter?: Filter<FilterAttribute>[];
}

interface TicketListResponse {
  tickets: TicketSummary[];

  // The total number of tickets the request would return without any limits
  total: number;
}
```

## Ticket Details

- **Endpoint:** `GET /tickets/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type TicketDetailsResponse = Ticket;
```

## Ticket status change

- **Endpoint:** `PATCH /tickets/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accesible if the user is an administrator

```ts
interface TicketPatchRequest {
  // Changes the status of the ticket to the passed-in value
  status?: TicketStatus;
}
```

If the status was changed successfully, simply return a 204 message (i.e. no content)

## Adding Ticket messages

- **Endpoint:** `POST /tickets/{id}/messages`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface TicketPostRequest {
  message: string;
}

type TicketPostResponse = TicketMessage;
```
