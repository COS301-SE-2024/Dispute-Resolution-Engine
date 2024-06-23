import { z } from "zod";

export const profileSchema = z.object({
  firstName: z.string().min(1, "Required"),
  surname: z.string().min(1, "Required"),
  phoneNumber: z.string().min(1, "Required"),

  gender: z.enum(["Male", "Female", "Non-binary", "Prefer not to say", "Other"]),
  country: z.string().min(1, "Required"),

  timezone: z.string().default("GMT+002"),
  preferredLanguage: z.string().default("en-US"),

  addrCountry: z.string().optional(),
  addrProvince: z.string().optional(),
  addrCity: z.string().optional(),

  addrStreet: z.string().optional(),
  addrStreet2: z.string().optional(),
  addrStreet3: z.string().optional(),
});

export type ProfileData = z.infer<typeof profileSchema>;
export type ProfileError = z.ZodFormattedError<ProfileData>;
