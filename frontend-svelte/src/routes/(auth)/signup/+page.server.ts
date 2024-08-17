import { signupSchema } from "$lib/schema/auth";
import { fail } from "@sveltejs/kit";
import { message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import type { Actions } from "./$types";
import { fetchCountries, fetchLanguages } from "$lib/api";
import { GENDERS } from "$lib/server";

export const load = async (event) => {
    const form = await superValidate(zod(signupSchema));
    const countries = (await fetchCountries(event.fetch)).data!;
    const languages = (await fetchLanguages()).data!;
    const genders = GENDERS.map(g => ({
        id: g,
        label: g
    }));

    return {
        form,
        countries,
        languages,
        genders
    };
};

export const actions = {
    default: async ({request}) => {
        const form = await superValidate(request, zod(signupSchema));
        if(!form.valid) {
            return fail(400, { form });
        }
        return message(form, "form posted successfully");
    }
} satisfies Actions;
