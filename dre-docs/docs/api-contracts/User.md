# User

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

# user analytics

- **Endpoint:** `POST /user/analytics`
- **Headers:**
  - `Autherization: Bearer <JWT>`

```ts
interface DateRange {
  startDate: string;
  endDate: string;
}

interface UserAnalytics {
  columnValueComparison?: Array<{ column: string; value: any }>;
  orderBy?: Array<{ column: string; direction: "asc" | "desc" }>;
  dateRanges?: {
    created_at?: DateRange;
    updated_at?: DateRange;
    last_login?: DateRange;
  };
  groupBy?: string[];
  count: boolean;
}
```
