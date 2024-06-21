import { AddressType, Gender } from ".";
import { Address } from "../schema/address";

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
  address: Address;
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

/**
 * Request used to reset the password of an authenticated user.
 *
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 */
export interface ResetPasswordRequest {
  oldPassword: string;
  newPassword: string;
}
/**
 * Returns a new updated JWT for the password reset
 */
export type ResetPasswordResponse = string;

export interface ForgotPasswordRequest {
  email: string;
}

/**
 * Returns a success message
 */
export type ForgotPasswordResponse = string;
