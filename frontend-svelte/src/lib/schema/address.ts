import { z } from "zod";

export const addressSchema = z.object({
  country: z.string().trim().min(1),
  province: z.string().trim().min(1),
  city: z.string().trim().min(1),

  address1: z.string().trim().min(1),
  address2: z.string().trim().min(1),
  address3: z.string().trim(),
});

export type Address = z.infer<typeof addressSchema>;
