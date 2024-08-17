export * from './auth';

/**
 * Generic Response from the API
 */
export type Result<T, E = string> =
	| {
			data: T;
			error?: never;
	  }
	| {
			data?: never;
			error: E;
	  };

export type FetchFn = (
	input: string | URL | globalThis.Request,
	init?: RequestInit
) => Promise<Response>;

export type Gender = 'Male' | 'Female' | 'Non-binary' | 'Prefer not to say' | 'Other';
export type AddressType = 'Postal' | 'Physical' | 'Billing';

export interface Address {
	address_type: AddressType;
	country: string;
	city: string;
	province: string;
	street: string;
	street2: string;
	street3: string;
}
