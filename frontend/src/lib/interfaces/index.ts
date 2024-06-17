export * from "./auth";

export type Result<T> =
  | {
      data: T;
      error?: never;
    }
  | {
      data?: never;
      error: string;
    };

export type Gender = "Male" | "Female" | "Non-binary" | "Prefer not to say" | "Other";
export type AddressType = "Postal" | "Physical" | "Billing";

export interface Address {
  address_type: AddressType;
  country: string;
  city: string;
  province: string;
  street: string;
  street2: string;
  street3: string;
}
