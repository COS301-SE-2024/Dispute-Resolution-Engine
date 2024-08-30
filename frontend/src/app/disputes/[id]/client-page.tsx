"use client";

import ExpertItem from "@/components/dispute/negotiator";
import FileInput from "@/components/form/file-input";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
  CardFooter,
} from "@/components/ui/card";
import { uploadEvidence } from "@/lib/actions/dispute";
import { DisputeResponse } from "@/lib/interfaces/dispute";
import { File as FileIcon } from "lucide-react";
import { FormEvent, useState } from "react";

export default function DisputeClientPage({ data }: { data: DisputeResponse }) {
  const [files, setFiles] = useState<File[]>([]);

  const onFilesSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    files.forEach((file) => formData.append("files", file, file.name));

    const res = await uploadEvidence(undefined, formData);
    console.log(res);
  };

  const resEvidence = data.evidence.filter((e) => e.uploader_role == "Respondent");
  const compEvidence = data.evidence.filter((e) => e.uploader_role == "Complainant");

  return (
    <>
      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Description</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-white/70 mt-4">{data.description}</p>
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          {(data.role ?? "") == "Complainant" ? (
            <>
              <CardTitle>Your Evidence</CardTitle>
              {compEvidence.length == 0 && (
                <CardDescription>You have not uploaded any evidence</CardDescription>
              )}
            </>
          ) : (
            <>
              <CardTitle>Complainant&apos;s Evidence</CardTitle>
              {compEvidence.length == 0 && (
                <CardDescription>The complainant has not uploaded any evidence</CardDescription>
              )}
            </>
          )}
        </CardHeader>
        <CardContent>
          <ul className="flex gap-2 flex-wrap">
            {compEvidence.map((evi, i) => (
              <li key={i} className="rounded-lg bg-gray-950 p-4 text-center text-gray-50 w-40">
                <a href={evi.url} title={evi.label}>
                  <FileIcon className="mx-auto" size="2rem" />
                  <p className="mt-2 text-sm font-medium truncate w-full">{evi.label}</p>
                </a>
              </li>
            ))}
          </ul>
          {(data.role ?? "") == "Complainant" && <UploadForm disputeId={data.id} />}
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          {(data.role ?? "") == "Respondent" ? (
            <>
              <CardTitle>Your Evidence</CardTitle>
              {resEvidence.length == 0 && (
                <CardDescription>You have not uploaded any evidence</CardDescription>
              )}
            </>
          ) : (
            <>
              <CardTitle>Respondant&apos;s Evidence</CardTitle>
              {resEvidence.length == 0 && (
                <CardDescription>The respondent has not uploaded any evidence</CardDescription>
              )}
            </>
          )}
        </CardHeader>
        <CardContent>
          <ul className="flex gap-2 flex-wrap">
            {resEvidence.map((evi, i) => (
              <li key={i} className="rounded-lg bg-gray-950 p-4 text-center text-gray-50 w-40">
                <a href={evi.url} title={evi.label}>
                  <FileIcon className="mx-auto" size="2rem" />
                  <p className="mt-2 text-sm font-medium truncate w-full">{evi.label}</p>
                </a>
              </li>
            ))}
          </ul>
          {(data.role ?? "") == "Respondent" && <UploadForm disputeId={data.id} />}
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Experts</CardTitle>
          {data.experts.length == 0 && (
            <CardDescription>No experts have been assigned yet.</CardDescription>
          )}
        </CardHeader>
        <CardContent>
          <ul>
            {data.experts.length > 0 &&
              data.experts.map((expert) => (
                <li key={expert.id}>
                  <ExpertItem {...expert} dispute_id={data.id} />
                </li>
              ))}
          </ul>
        </CardContent>
      </Card>
    </>
  );
}

function UploadForm({ disputeId }: { disputeId: string }) {
  const [files, setFiles] = useState<File[]>([]);

  const onFilesSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    files.forEach((file) => formData.append("files", file, file.name));

    const res = await uploadEvidence(undefined, formData);
    console.log(res);
  };
  return (
    <form onSubmit={onFilesSubmit} className="space-y-4">
      <input type="hidden" name="dispute_id" value={disputeId} />
      <FileInput onValueChange={setFiles} />
      <Button disabled={files.length == 0} type="submit">
        Upload
      </Button>
    </form>
  );
}
