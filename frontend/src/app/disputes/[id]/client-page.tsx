"use client";

import FileInput from "@/components/form/file-input";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { uploadEvidence } from "@/lib/actions/dispute";
import { DisputeResponse } from "@/lib/interfaces/dispute";
import { File as FileIcon } from "lucide-react";
import { ChangeEvent, FormEvent, ReactNode, useState } from "react";

export default function DisputeClientPage({ data }: { data: DisputeResponse }) {
  const [files, setFiles] = useState<File[]>([]);

  const onFilesSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    files.forEach((file) => formData.append("files", file, file.name));

    const res = await uploadEvidence(undefined, formData);
    console.log(res);
  };

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
            <CardTitle>Your Evidence</CardTitle>
          ) : (
            <CardTitle>Complainant&apos;s Evidence</CardTitle>
          )}
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent>
          <ul className="flex gap-2 flex-wrap">
            {data.evidence.map((evi, i) => (
              <li key={i} className="rounded-lg bg-gray-950 p-4 text-center text-gray-50 w-40">
                <a href={evi.url} title={evi.label}>
                  <FileIcon className="mx-auto" size="2rem" />
                  <p className="mt-2 text-sm font-medium truncate w-full">{evi.label}</p>
                </a>
              </li>
            ))}
          </ul>
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          {(data.role ?? "") == "Respondant" ? (
            <CardTitle>Your Evidence</CardTitle>
          ) : (
            <CardTitle>Respondant&apos;s Evidence</CardTitle>
          )}
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent>
          <ul className="flex gap-2 flex-wrap">
            {data.evidence.map((evi, i) => (
              <li key={i} className="rounded-lg bg-gray-950 p-4 text-center text-gray-50 w-40">
                <a href={evi.url} title={evi.label}>
                  <FileIcon className="mx-auto" size="2rem" />
                  <p className="mt-2 text-sm font-medium truncate w-full">{evi.label}</p>
                </a>
              </li>
            ))}
          </ul>
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onFilesSubmit}>
            <input type="hidden" name="dispute_id" value={data.id} />
            <FileInput onValueChange={setFiles} />
            <Button disabled={files.length == 0} type="submit">
              Upload
            </Button>
          </form>
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Experts</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onFilesSubmit}>
            <input type="hidden" name="dispute_id" value={data.id} />
            <FileInput onValueChange={setFiles} />
            <Button disabled={files.length == 0} type="submit">
              Upload
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  );
}
