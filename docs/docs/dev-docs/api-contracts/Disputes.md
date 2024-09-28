# Disputes

# Utility Types

```ts
type DisputeSummary = {
  id: string;
  title: string;
  description: string;
  status: string;
  role: "Complainant" | "Respondant";
};

type Evidence = {
  label: string;
  url: string;
  date_submitted: string;
};

type Expert = {
  id: string;
  full_name: string;
  email: string;
  phone: string;
  role: string;
};
```

# Utility Functions

- **Endpoint:** `GET /utils/dispute_statuses`
- **Headers:**
  - None expected
- Will return a list of all possible states a dispute can be in:

```ts
type DisputeStatusesResponse = string[];
```

# Dispute Summaries

- **Endpoint:** `GET /disputes`
- **Headers:**
  - `Authorization: Bearer <JWT>`

Will return a list of dispute summaries the user is involved in:

```ts
type DisputeListResponse = DisputeSummary[];
```

# Dispute Details

- **Endpoint:** `GET /disputes/{id}`
- **Headers:**
  - `Authorization: Bearer <JWT>`

Will return detailed information about a dispute the user is involved in:

```ts
type DisputeResponse = {
  id: string;
  title: string;
  description: string;
  status: string;
  date_created: string;

  evidence: Evidence[];
  experts: Expert[];
};
```

# Dispute Creation

- **Endpoint:** `POST /disputes/create`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type DisputeCreateRequest = {
  title: string;
  description: string;
  "respondent[full_name]": string;
  "respondent[email]": string;
  "respondent[telephone]": string;
  "respondent[workflow]": string;
  files: File;
};
```

The response will return the ID of the newly-created dispute:

```ts
type DisputeCreateResponse = {
  id: number;
};
```

# Dispute Evidence upload

- **Endpoint:** `POST /disputes/{id}/evidence`
- **Headers:**
  - `Authorization: Bearer <JWT>`

**NOTE:** This endpoint involves files being exchanged. This requires the use of `multipart/form-data` instead of JSON.
This is left incomplete until this exception is documenteda.

```ts
interface EvidenceUploadRequest {
  file: File[];
}
```

# Dispute Status Change

- **Endpoint:** `PUT /disputes/{id}/status`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
type UpdateRequest = {
  status: string;
};
```
