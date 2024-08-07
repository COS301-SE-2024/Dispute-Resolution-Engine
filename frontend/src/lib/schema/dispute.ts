import { z } from "zod";

export const disputeCreateSchema = z.object({
  title: z.string().min(2).max(50),
  respondentName: z.string().min(1).max(50),
  respondentEmail: z.string().email(),
  respondentTelephone: z.string().min(10).max(15),
  summary: z.string().min(3).max(500),

  // Dummy variable to make RHF happy
  file: z.any(),
});

export type DisputeCreateData = z.infer<typeof disputeCreateSchema>;
export type DisputeCreateError = z.ZodFormattedError<DisputeCreateData>;

export const expertRejectSchema = z.object({
  dispute_id: z.string().min(1),
  expert_id: z.string().min(1),
  reason: z.string().min(20),
});
export type ExpertRejectData = z.infer<typeof expertRejectSchema>;
export type ExpertRejectError = z.ZodFormattedError<ExpertRejectData>;
