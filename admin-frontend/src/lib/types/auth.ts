import { JwtHeader } from "jwt-decode";

export type UserStatus = "Active" | "Inactive" | "Suspended";
export type Gender = "Male" | "Female" | "Non-binary" | "Prefer not to say" | "Other";

export interface User {
  id: number;
  first_name?: string;
  surname?: string;
  birthdate?: string;
  nationality?: string;
  role?: string;
  email: string;
  phone_number?: string;
  address_id?: number;
  status?: UserStatus;
  gender?: string;
  preferred_language?: string;
  timezone?: string;
}

export type UserJwt = JwtHeader & {
  user: User;
};
