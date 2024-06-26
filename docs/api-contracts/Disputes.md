All endpoints follow the following general type:
```ts
type Result<T> =
  | {
      data: T;
      error?: never;
    }
  | {
      data?: never;
      error: string;
    };
```

Which corresponds to either returning:
```json5
{
    "data": /* ... some data */
}
```
or
```json5
{
    "error": "error message"
}
```

# Utility Types
```ts
type DisputeSummary = {
  id: string;
  title: string;
  description: string;
  status: string;
};

type Evidence = {
  label: string;
  url: string;
  date_submitted: string;
};

type Expert = {
  full_name: string;
  email: string;
  phone: string;
};
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

**NOTE:** This endpoint involves files being exchanged. This requires the use of `multipart/form-data` instead of JSON.
Below is temporary description until this exception is documented

```ts
interface DisputeCreateRequest {
  title: string;
  description: string;
  evidence: {
    label: string;
    file: File;
  }[];
  respondent: {
    full_name: string;
    email: string;
    telephone: string;
  };
}

```

The response will return a success message
```ts
type DisputeCreateResponse = string;
```

# Dispute Evidence upload
- **Endpoint:** `POST /disputes/{id}/upload`
- **Headers:**
    - `Authorization: Bearer <JWT>`

**NOTE:** This endpoint involves files being exchanged. This requires the use of `multipart/form-data` instead of JSON.
This is left incomplete until this exception is documenteda.