import { z } from "zod";

export const timerSchema = z.object({
  duration: z.string(),
  on_expire: z.string(),
});

export const eventSchema = z.object({
  label: z.string(),
  next_state: z.string(),
});

export const stateSchema = z.object({
  label: z.string(),
  description: z.string(),
  events: z.record(eventSchema),
  timer: timerSchema.optional(),
});

export const workflowSchema = z.object({
  label: z.string(),
  initial: z.string(),
  states: z.record(stateSchema),
});
