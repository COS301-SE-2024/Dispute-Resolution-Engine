"use client";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { DisputeCreateRequest } from "@/lib/interfaces/dispute";
import { DisputeCreateData, disputeCreateSchema } from "@/lib/schema/dispute";
import { useRef } from "react";
import { createDispute } from "@/lib/actions/dispute";

export default function CreateDisputeClient() {
  const form = useForm<DisputeCreateData>({
    defaultValues: {
      title: "",
      respondentName: "",
      respondentEmail: "",
      respondentTelephone: "",
      summary: "",
    },
    resolver: zodResolver(disputeCreateSchema),
  });

  const formRef = useRef(null);

  const onSubmit = async function (dataFromForm: DisputeCreateData) {
    createDispute(null, new FormData(formRef.current!));
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
          Create a Dispute
        </CardTitle>
      </CardHeader>
      <Form {...form}>
        <form ref={formRef} onSubmit={form.handleSubmit(onSubmit)} className="w-full pt-0 p-10">
          <div className="space-y-5">
            <FormField
              control={form.control}
              name="title"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Title</FormLabel>
                  <FormControl>
                    <Input placeholder="Dispute Title" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="respondentName"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Respondent Name</FormLabel>
                  <FormControl>
                    <Input placeholder="John Doe" {...field} />
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
                    <Input placeholder="012 345 6789" {...field} />
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
                    <Input placeholder="The aforementioned party..." {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="file"
              render={({ field }) => {
                return (
                  <FormItem>
                    <FormLabel>Evidence</FormLabel>
                    <FormControl>
                      <Input type="file" placeholder="shadcn" name="file" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                );
              }}
            />
            <Button type="submit">Submit</Button>
          </div>
        </form>
      </Form>
    </Card>
  );
}
