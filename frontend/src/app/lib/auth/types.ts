import { z } from "zod";

export const signupSchema = z
  .object({
    firstName: z.string().min(1, "Required"),
    lastName: z.string().min(1, "Required"),
    email: z.string().min(1, "Required").email("Please enter a valid email"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters long")
      .regex(/\d/gm, "Password must contain at least one digit")
      .regex(/[A-Za-z]/gm, "Password must contain at least one letter")
      .regex(/[^\w\d\s:]/gm, "Password must contain a special character"),
    passwordConfirm: z.string(),

    addrCountry: z.string().min(1, "Required"),
    addrProvince: z.string().min(1, "Required"),
    addrCity: z.string().min(1, "Required"),
    addrStreet3: z.string().min(1, "Required"),
    addrStreet2: z.string().min(1, "Required"),
    addrStreet: z.string().min(1, "Required"),

    dateOfBirth: z.string().date("Invalid date"),
    idNumber: z.string().regex(/\d{13}/gm, "Invalid ID number"),
  })
  .superRefine((arg, ctx) => {
    if (arg.password !== arg.passwordConfirm) {
      ctx.addIssue({
        code: "custom",
        message: "The passwords did not match",
        path: ["passwordConfirm"],
      });
    }
  });

export type SignupData = z.infer<typeof signupSchema>;
export type SignupError = z.ZodFormattedError<SignupData>;
