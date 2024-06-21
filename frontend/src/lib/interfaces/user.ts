import { Address, Gender } from ".";

/**
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 */
export interface UserProfileResponse {
  first_name: string;
  surname: string;
  email: string;
  phone_number: string;

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

  addresses: Address[];
  useDarkMode: boolean;
}

/**
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 */
export interface UserProfileUpdateRequest {
  first_name: string;
  surname: string;
  phone_number: string;
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
  addresses: Address[];
}

/**
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 */
export interface UserProfileUpdateResponse extends UserProfileResponse {}

/**
 * Add a HTTP header with authentication header:
 * ```http
 * Authorization: Bearer <JWT>
 * ```
 */
export type UserProfileRemoveResponse = string;
