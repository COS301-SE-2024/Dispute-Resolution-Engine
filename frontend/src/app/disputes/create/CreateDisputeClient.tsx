"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { FormProvider, useForm } from "react-hook-form";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { DisputeCreateData, disputeCreateSchema } from "@/lib/schema/dispute";
import { useId, useRef, useState } from "react";
import { createDispute } from "@/lib/actions/dispute";
import { Textarea } from "@/components/ui/textarea";
import FileInput from "@/components/form/file-input";
import { FormField, FormMessage } from "@/components/ui/form-client";

const CreateMessage = FormMessage<DisputeCreateData>;
const CreateField = FormField<DisputeCreateData>;

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

  const { register, setError } = form;

  const [files, setFiles] = useState<File[]>([]);

  // Used to access the FormData on a form submission
  const formRef = useRef(null);
  const onSubmit = async function (dataFromForm: DisputeCreateData) {
    const formdata = new FormData(formRef.current!);
    files.forEach((file) => formdata.append("file", file, file.name));

    const res = await createDispute(null, formdata);
    if (res && res.error) {
      setError("root", { type: "custom", message: res.error });
    }
  };

  const resNameId = useId();
  const resEmailId = useId();
  const resTelId = useId();

  const titleId = useId();
  const summaryId = useId();
  const fileId = useId();

  return (
    <FormProvider {...form}>
      <form ref={formRef} onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-10 py-5">
        <Card>
          <CardHeader>
            <CardTitle className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
              Respondant Information
            </CardTitle>
            <CardDescription>Who are you filing a dispute against?</CardDescription>
          </CardHeader>
          <CardContent className="space-y-5">
            <CreateField id={resNameId} name="respondentName" label="Respondent Name">
              <Input id={resNameId} placeholder="John Doe" {...register("respondentName")} />
            </CreateField>
            <CreateField id={resEmailId} name="respondentEmail" label="Respondent Email">
              <Input
                id={resEmailId}
                placeholder="abc@example.com"
                {...register("respondentEmail")}
              />
            </CreateField>
            <CreateField id={resTelId} name="respondentTelephone" label="Respondent Telephone">
              <Input
                id={resTelId}
                placeholder="012 345 6789"
                {...register("respondentTelephone")}
              />
            </CreateField>
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
            <CreateField id={titleId} name="title" label="Dispute Title">
              <Input id={titleId} placeholder="Dispute Title" {...register("title")} />
            </CreateField>
            <CreateField id={summaryId} name="summary" label="Dispute Summary">
              <Textarea
                id={summaryId}
                placeholder="A short description of why you are filing a dispute (max. 500 words)"
                {...register("summary")}
              />
            </CreateField>
            <CreateField id={fileId} name="file" label="Evidence">
              <FileInput id={fileId} onValueChange={setFiles} />
            </CreateField>
          </CardContent>
          <CardFooter className="flex justify-end items-center gap-2">
            <CreateMessage />
            <Button type="submit">Create</Button>
          </CardFooter>
        </Card>
      </form>
    </FormProvider>
  );
}
