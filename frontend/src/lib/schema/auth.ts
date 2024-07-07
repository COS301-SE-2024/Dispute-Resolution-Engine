import { z } from "zod";

export const signupSchema = z
  .object({
    firstName: z.string().min(1, "Required"),
    lastName: z.string().min(1, "Required"),
    email: z.string().min(1, "Required").email("Please enter a valid email"),
    gender: z.string().min(1, "Required"),
    nationality: z.string().min(1, "Required"),
    preferredLanguage: z.string().min(1, "Required"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters long")
      .regex(/\d/gm, "Password must contain at least one digit")
      .regex(/[A-Za-z]/gm, "Password must contain at least one letter")
      .regex(/[^\w\d\s:]/gm, "Password must contain a special character"),
    passwordConfirm: z.string(),

    dateOfBirth: z.string().date("Invalid date"),
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

export const loginSchema = z.object({
  email: z.string().min(1, "Required"),
  password: z.string().min(1, "Required"),
});
export type LoginData = z.infer<typeof loginSchema>;
export type LoginError = z.ZodFormattedError<LoginData>;

export const verifySchema = z.object({
  pin: z.string().length(6),
});
export type VerifyData = z.infer<typeof verifySchema>;
export type VerifyError = z.ZodFormattedError<VerifyData>;
