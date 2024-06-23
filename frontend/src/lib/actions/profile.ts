"use server";

import { UserAddressUpdateRequest, UserProfileUpdateResponse } from "../interfaces/user";
import { ProfileError, profileSchema } from "../schema/profile";
import { Result } from "../types";
import { getAuthToken } from "../util/jwt";
import { API_URL } from "../utils";

async function updateAddress(req: UserAddressUpdateRequest): Promise<Result<string>> {
    const response = await fetch(`${API_URL}/user/profile/address`, {
        method: "PUT",
        headers: {
            Authorization: `Bearer ${getAuthToken()}`,
        },
        body: JSON.stringify(req),
    });
        return response.clone().json()
        .catch(async (e: Error) => ({
            error: e.message + ": " + await response.text(),
        }));
}

export async function updateProfile(
    _init: unknown,
    formData: FormData
): Promise<Result<UserProfileUpdateResponse, ProfileError>> {
    const { data: parsed, error: parseErr } = profileSchema.safeParse(Object.fromEntries(formData));
    if (parseErr) {
        console.error(parseErr.format());
        return {
            error: parseErr.format(),
        };
    }

    const updateRes = await fetch(`${API_URL}/user/profile`, {
        method: "PUT",
        headers: {
            Authorization: `Bearer ${getAuthToken()}`,
        },
        body: JSON.stringify({
            first_name: parsed.firstName,
            surname: parsed.surname,
            phone_number: parsed.phoneNumber,

            gender: parsed.gender,
            nationality: parsed.country,

            timezone: parsed.timezone,
            preferred_language: parsed.preferredLanguage
        }),
    });
    const { data, error } = await updateRes.clone().json()
        .catch(async (e: Error) => ({
            error: e.message + ": " + (await updateRes.text()) ,
        }));

    if (error) {
        return {
            error: {
                _errors: [error],
            },
        };
    }

    const res = await updateAddress({
        country: parsed.addrCountry,
        province: parsed.addrProvince,
        city: parsed.addrCity,
        street3: parsed.addrStreet3,
        street2: parsed.addrStreet2,
        street: parsed.addrStreet,
        address_type: "Postal",
    });
    if (res.error) {
        return {
            error: {
                _errors: [res.error],
            },
        };
    }
    return { data };
}


export async function deleteProfile(
    _init: unknown,
    _formData: FormData
): Promise<Result<string>> {
    return fetch(`${API_URL}/user/profile`, {
        method: "DELETE",
        headers: {
            Authorization: getAuthToken()
        }
    }).then(res => res.json()).catch((e: Error) => ({
        error: e.message
    }))
}
