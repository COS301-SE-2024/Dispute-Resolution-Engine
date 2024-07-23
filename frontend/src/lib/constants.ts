export const JWT_VERIFY_TIMEOUT = 120;
export const JWT_TIMEOUT = 60 * 20;
export const JWT_KEY = "jwt";

// Type definitino needed to make Zod happy with z.enum
export const GENDERS: [string, ...string[]] = [
  "Male",
  "Female",
  "Non-binary",
  "Prefer not to say",
  "Other",
];
