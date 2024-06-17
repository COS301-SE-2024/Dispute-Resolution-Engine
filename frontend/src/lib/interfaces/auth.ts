export type Gender = "Male" | "Female" | "Non-binary" | "Prefer not to say" | "Other";
export type AddressType = "Postal" | "Physical" | "Billing";

export interface LoginRequest {
  email: string;
  password: string;
}

/**
 * Returns the JWT of the logged-in user
 */
export type LoginResponse = string;

export interface SignupRequest {
  first_name: string;
  surname: string;
  email: string;
  phone_number: string;

  password: string;

  birthdate: string;
  gender: Gender;
  nationality: string;

  /**
   * Hard-coded for now (will be overwritten by API)
   */
  timezone: string;

  /**
   * Hard-coded for now (will be overwritten by API)
   */
  preferred_language: string;

  // Address information
  address_type: AddressType;
  country: string;
  city: string;
  province: string;
  street: string;
  street2: string;
  street3: string;
}

/**
 * Returns a temporary JWT (pending email verification)
 */
export type SignupResponse = string;

/**
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 */
export interface VerifyRequest {
  pin: string;
}

/**
 * Returns the JWT of the verified user
 */
export type VerifyResponse = string;
