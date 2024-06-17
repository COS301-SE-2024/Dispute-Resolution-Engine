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

  // Stub for the time being
  timezone: string;
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
 * Simply returns a success message (so no extra information)
 */
export type SignupResponse = string;
