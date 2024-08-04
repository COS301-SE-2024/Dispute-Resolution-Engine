# Authentication

# Login

- **Endpoint:** `POST /auth/login`

```ts
interface LoginRequest {
  email: string;
  password: string;
}
```

The response returns the JWT of the logged-in user:

```ts
type LoginResponse = string;
```

# Signup

- **Endpoint:** `POST /auth/signup`

```ts
interface SignupRequest {
  first_name: string;
  surname: string;
  email: string;
  phone_number: string;

  password: string;

  birthdate: string;
  gender: Gender;
  nationality: string;

  timezone: string;
  preferred_language: string; // en-US by default

  // Address information
  // address: Address;
}
```

The response will return a temporary JWT that can be used for email verification:

```ts
type SignupResponse = string;
```

# Email Verification

- **Endpoint:** `POST /auth/signup`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface VerifyRequest {
  pin: string;
}
```

The response will invalidate the existing JWT and return a new JWT for the full verified user:

```ts
type VerifyResponse = string;
```

# Password Resetting

- **Endpoint:** `POST /auth/reset-password/send-email`

```ts
interface SendResetRequest {
  email: string;
}
```

The response will return a success message

```ts
type SendResetResponse = string;
```

## Password Resetting (Resetting)

- **Endpoint:** `POST /auth/reset-password/reset`
- **Headers:**
  - `Authorization: Bearer <Temp-JWT>`

```ts
interface ResetPasswordRequest {
  newPassword: string;
}
```

The response will return a success message

```ts
type ResetPasswordResponse = string;
```
