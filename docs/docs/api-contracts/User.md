# User

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
  error: "error message",
}
```

# User Profile

- **Endpoint:** `GET /user/profile`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface UserProfileResponse {
  first_name: string;
  surname: string;
  email: string;
  phone_number: string;

  birthdate: string;
  gender: Gender;
  nationality: string;

  timezone: string;
  preferred_language: string;

  addresses: Address[];
  theme: "light" | "dark";
}
```

# Updating User Profile

- **Endpoint:** `PUT /user/profile`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface UserProfileUpdateRequest {
  first_name: string;
  surname: string;
  phone_number: string;

  gender: Gender;
  nationality: string;

  timezone: string;
  preferred_language: string;

  addresses: Address[];
  theme: "light" | "dark";
}
```

The server will respond with the new updated user information:

```ts
type UserProfileUpdateResponse = string;
```

# Deleting Account

- **Endpoint:** `DELETE /user/profile`
- **Headers:**
  - `Authorization: Bearer <JWT>`

The response will simply be a success message

```ts
type UserProfileRemoveResponse = string;
```

# Update Address

- **Endpoint:** `PUT /user/profile/address`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface UserAddressUpdateRequest {
  country?: string;
  province?: string;
  city?: string;
  street3?: string;
  street2?: string;
  street?: string;
  address_type?: string;
}
```

The server will respond with a success or failure

```ts
type UserAddressUpdateResponse = string;
```

<<<<<<< HEAD:docs/docs/api-contracts/User.md

# user analytics

=======

# User analytics

> > > > > > > feat/analytics-ui:docs/api-contracts/User.md

- **Endpoint:** `POST /user/analytics`
- **Headers:**
  - `Autherization: Bearer <JWT>`

```ts
interface DateRange {
  startDate: string;
  endDate: string;
}

<<<<<<< HEAD:docs/docs/api-contracts/User.md
interface UserAnalytics {
  columnValueComparison?: Array<{ column: string; value: any }>;
  orderBy?: Array<{ column: string; direction: "asc" | "desc" }>;
=======
interface UserAnalyticsRequest {
  // Used to fuzzy search columns from the users table, e.g. { "column": "role", value: "Mediator" }
  columnValueComparison?: Array<{ column: string; value: any }>;

  // Used to order the results obtained. Multiple elements are only useful when a column contains items with the same key
  orderBy?: Array<{ column: string; direction: "asc" | "desc" }>;

>>>>>>> feat/analytics-ui:docs/api-contracts/User.md
  dateRanges?: {
    created_at?: DateRange;
    updated_at?: DateRange;
    last_login?: DateRange;
  };

  // Column names to group by
  groupBy?: string[];

  // Whether to return a count instead of concrete results
  count: boolean;
}

type UserAnalyticsResponse =
  | {
      count: number;
    }
  | {
      users: Array<{
        id: number;
        first_name: string;
        surname: string;
        birthdate: string;
        nationality: string;
        role: string;
        email: string;
        phone_number?: string;
        address_id?: number;
        created_at: string;
        updated_at?: string;
        last_login?: string;
        status: string;
        gender: string;
        preferred_language?: string;
        timezone?: string;
      }>;
    };
```
