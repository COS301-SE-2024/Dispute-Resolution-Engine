import type { FetchFn, Result } from '$lib/interfaces';
import { API_URL } from '$lib/server';

export type Country = {
	id: number;
	country_code: string;
	country_name: string;
};

export type Language = {
	id: string;
	label: string;
};

export async function fetchCountries(fetch: FetchFn): Promise<Result<Country[]>> {
	return fetch(`${API_URL}/utils/countries`, {}).then((r) => r.json());
}

export async function fetchLanguages(): Promise<Result<Language[]>> {
	return {
		data: [
			{
				id: 'en-US',
				label: 'English'
			}
		]
	};
}
