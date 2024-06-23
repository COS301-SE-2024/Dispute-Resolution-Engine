"use client"

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { DisputeCreateRequest } from "@/lib/interfaces/dispute";
import { createDispute } from "@/lib/api/dispute";

const formSchema = z.object({
  title: z.string().min(2).max(50),
  respondentName: z.string().min(1).max(50),
  respondentEmail: z.string().email(),
  respondentTelephone: z.string().min(10).max(15),
  summary: z.string().min(3).max(500),
  file: z.instanceof(FileList).optional()
});

export default function CreateDisputeClient() {
  const form = useForm<z.infer<typeof formSchema>>({
    defaultValues: {
      title: "",
      respondentName: "",
      respondentEmail: "",
      respondentTelephone: "",
      summary: ""
    },
    resolver: zodResolver(formSchema)
  });

  const fileRef = form.register("file");

  const onSubmit = async function(dataFromForm: z.infer<typeof formSchema>) {
    if (dataFromForm.file === undefined) {
      dataFromForm.file = new FileList();
    }

    const requestData: DisputeCreateRequest = {
      title: dataFromForm.title,
      description: dataFromForm.summary,
      evidence: [...dataFromForm.file],
      respondent: {
        full_name: dataFromForm.respondentName,
        email: dataFromForm.respondentEmail,
        telephone: dataFromForm.respondentTelephone
      }
    };

// Create a new FormData object
    const formData = new FormData();

// Append each property of requestData to formData
    formData.append("title", requestData.title);
    formData.append("description", requestData.description);
    requestData.evidence.forEach((file, index) => {
      formData.append(`evidence[${index}]`, file);
    });
    formData.append("respondent[full_name]", requestData.respondent.full_name);
    formData.append("respondent[email]", requestData.respondent.email);
    formData.append("respondent[telephone]", requestData.respondent.telephone);
    const response = await createDispute(formData);
    console.log(response);
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">Create a
          Dispute</CardTitle>
      </CardHeader>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="w-full pt-0 p-10">
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
                      <Input type="file" placeholder="shadcn" {...fileRef} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                );
              }}
            />
            <Button type="submit">Submit</Button>
          </div>
        </form>
      </Form></Card>
  );
}