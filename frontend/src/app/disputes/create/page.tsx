"use client";

import { z } from "zod";
import { useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { API_URL } from "@/lib/utils";
import { Evidence } from "@/lib/interfaces/dispute";
export interface DisputeCreateRequest {
  description: string;
  desired_outcome: string;

  respondent: {
    full_name: string;
    email: string;
    telephone: string;
  };

  jurisdictional_basis: Evidence;

  /**
   * IDs of all adjudicators to be appointed
   */
  adjudicators: string[];

  // This should be FormData, but I don't know how to annotate that
  evidence: Evidence[];
}
// Add field for evidence which is a file upload
const formSchema = z.object({
  title: z.string().min(2).max(50),
  respondentEmail: z.string().email(),
  respondentTelephone: z.string().min(10).max(15),
  summary: z.string().max(500),
})

export default function CreateDispute() {
  const form = useForm<z.infer<typeof formSchema>>({
    defaultValues: {
      title: "",
      respondentEmail: "",
      respondentTelephone: "",
      summary: "",
    },
    resolver: zodResolver(formSchema),
  })
  function onSubmit(values: z.infer<typeof formSchema>) {
    const data = fetch(`${API_URL}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(values),
    })
    console.log(values)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Title</FormLabel>
              <FormControl>
                <Input placeholder="We be beefing" {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="respondentEmail"
          render={({ field }) => (
            <FormItem>
              <FormLabel>RespondentEmail</FormLabel>
              <FormControl>
                <Input placeholder="abc@example.com" {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="respondentTelephone"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Respondent Telephone</FormLabel>
              <FormControl>
                <Input placeholder="Mr Biggest Op" {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="summary"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Summary</FormLabel>
              <FormControl>
                <Input placeholder="He stole my chib" {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}
