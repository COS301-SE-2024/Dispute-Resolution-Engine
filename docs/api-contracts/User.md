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
- **Endpoint:** `PATCH /user/profile`
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
interface UserProfileUpdateResponse extends UserProfileResponse {}
```

# Deleting Account
- **Endpoint:** `DELETE /user/profile`
- **Headers:**
    - `Authorization: Bearer <JWT>`

The response will simply be a success message
```ts
type UserProfileRemoveResponse = string;
```