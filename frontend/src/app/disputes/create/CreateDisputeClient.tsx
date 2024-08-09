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
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { DisputeCreateData, disputeCreateSchema } from "@/lib/schema/dispute";
import { useId, useRef, useState } from "react";
import { createDispute } from "@/lib/actions/dispute";
import { Textarea } from "@/components/ui/textarea";
import FileInput from "@/components/form/file-input";
import { string } from "zod";

export default function CreateDisputeClient() {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
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
  const filesId = useId();

  const [files, setFiles] = useState<File[]>([]);

  const onSubmit = async function (dataFromForm: DisputeCreateData) {
    const formdata = new FormData(formRef.current!);
    files.forEach((file) => formdata.append("file", file, file.name));

    const res = createDispute(null, formdata);
    const comRes = (await res)
    if (comRes.error) {
      const data = comRes.error;
      setErrorMessage(comRes.error._errors[0])
    }
  };

  return (
    <Form {...form}>
      <form ref={formRef} onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-10 py-5">
        <Card>
          <CardHeader>
            <CardTitle className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
              Respondant Information
            </CardTitle>
            <CardDescription>Who are you filing a dispute against?</CardDescription>
          </CardHeader>
          <CardContent className="space-y-5">
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
          </CardContent>
        </Card>
        <Card className="w-full">
          <CardHeader>
            <CardTitle className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
              Dispute Details
            </CardTitle>
            <CardDescription>What is the dispute about?</CardDescription>
          </CardHeader>
          <CardContent className="space-y-5">
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
              name="summary"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Summary</FormLabel>
                  <FormControl>
                    <Textarea placeholder="The aforementioned party..." {...field} />
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
                      <FileInput onValueChange={setFiles} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                );
              }}
            />
            <Button type="submit">Create</Button>
          </CardContent>
        </Card>
      </form>
    </Form>
  );
}
