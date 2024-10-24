import { z } from "zod";

export const createTicketSchema = z.object({
  dispute_id: z.coerce.number({
    message: "Invalid dispute ID",
  }),
  subject: z.string().trim().min(1, "Subject is required"),
  body: z.string().trim().min(10, "Body must be a at least 10 characters"),
});

export type CreateTicketData = z.infer<typeof createTicketSchema>;
export type CreateTicketErrors = z.ZodFormattedError<CreateTicketData>;
