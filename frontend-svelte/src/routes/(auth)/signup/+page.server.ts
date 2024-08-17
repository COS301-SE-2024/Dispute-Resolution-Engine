import { signupSchema } from "$lib/schema/auth";
import { fail } from "@sveltejs/kit";
import { message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import type { Actions } from "./$types";

export const load = async () => {
    const form = await superValidate(zod(signupSchema));
    return { form };
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
