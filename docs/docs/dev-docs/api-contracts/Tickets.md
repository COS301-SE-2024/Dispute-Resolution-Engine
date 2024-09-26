# Tickets

The following types are used to describe tickets:

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

## Ticket List (Admin)

- **Endpoint:** `POST /tickets`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type SortAttribute =
  | "date_created" // The date is was created
  | "subject" // The subject of the ticket
  | "status" // The status of the ticket
  | "user"; // The full name of the user;

type TicketStatus = 
  'Open'   |
  'Closed' |
  'Solved' |
  'On Hold'


type FilterAttribute = {
  "attr": "status"
  "value": TicketStatus
}

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

## Ticket List (User)
- **Endpoint:** `POST /tickets`
- **Headers**
  - `Authorization: Bearer <JWT>`


```ts
Request = TicketListRequest;
Response = TicketListResponse;
```

## Ticket Details

- **Endpoint:** `GET /tickets/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** A normal user should not be able to access tickets that they did not create

```ts
type TicketDetailsResponse = Ticket;
```

## Ticket status change (Admin)

- **Endpoint:** `PATCH /tickets/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should only be accesible to administrators

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
- **Note:** Should only be accesible to administrators or the user that created the ticket

```ts
interface TicketPostRequest {
  message: string;
}

type TicketPostResponse = TicketMessage;
```

## Creating a ticket

- **Endpoint:** `POST /dispute/{id}/tickets`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Notes:**
  - Admins cannot create tickets, as they cannot be involved in disputes
  - Users should not be able to create tickets on disputes that they are not involved in

```ts
interface TicketCreateRequest {
  // The subject of the ticket
  subject: string;

  // The body of the ticket
  body: string;
}

type TicketCreateResponse = Omit<Ticket, "user">;
```
