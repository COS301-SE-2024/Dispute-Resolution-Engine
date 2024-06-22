import { z } from "zod";
import { addressSchema } from "./address";

export const profileSchema = z.object({
  first_name: z.string().min(1, "Required"),
  surname: z.string().min(1, "Required"),
  phone_number: z.string().min(1, "Required"),
  gender: z.enum(["Male", "Female", "Non-binary", "Prefer not to say", "Other"]),
  nationality: z.string().min(1, "Required"),

  timezone: z.string(),
  preferred_language: z.string().min(1, "Required"),
  addresses: z.array(addressSchema),
});

export type ProfileData = z.infer<typeof profileSchema>;
export type ProfileError = z.ZodFormattedError<ProfileData>;
