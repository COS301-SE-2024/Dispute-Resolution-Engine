# Experts

TODO

# Expert Objections

```ts
type ObjectionStatus = "Review" | "Sustained" | "Overruled";
export interface Objection {
  // ID of the objection itself.
  id: number;
  // ID  of the ticket the objection is related to
  ticket_id: number;

  // The full name of the expert being objected to
  expert_name: string;

  // The full name of the user that submitted the objection
  user_name: string;

  // When the objection was submitted
  date_submitted: string;

  // The status of the objection
  status: ObjectionStatus;
}
```

## Creating an objection

- **Endpoint:** `POST /disputes/{id}/objections`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Administrators should not be able to create objections

```ts
interface CreateObjectionRequest {
  expert_id: number;
  reason: string;
}
```

Return the created ticket ID:

```ts
export type CreateObjectionResponse = number;
```

## Retrieving objections

- **Endpoint:** `POST /disputes/experts/objections`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** Should not be accessible by an administrator
```ts
interface ViewExpertRejectionsRequest {
	expert_id?  : number
	dispute_id? : number
	limits?     : number   
	offset?     : number 
}
```

returns:
```ts
type ObjectionListResponse = Objection[];
```

## Reviewing an objection

- **Endpoint:** `PATCH /disputes/objections/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface ObjectionListRequest {
  status?: ObjectionStatus;
}
```

The response will be a 204 (no content) code.
