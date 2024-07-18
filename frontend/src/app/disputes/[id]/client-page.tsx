"use client";

import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { uploadEvidence } from "@/lib/actions/dispute";
import { DisputeResponse } from "@/lib/interfaces/dispute";
import { File as FileIcon } from "lucide-react";
import { ChangeEvent, FormEvent, FormEventHandler, useState } from "react";

export default function DisputeClientPage({ data }: { data: DisputeResponse }) {
  const [files, setFiles] = useState<File[]>([]);
  const onFilesChange = async (ev: ChangeEvent<HTMLInputElement>) => {
    setFiles([...files, ...ev.target.files!]);
    ev.target.value = "";
  };
  const removeFile = async (i: number) => {
    setFiles(files.filter((_f, j) => i !== j));
  };

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
          <CardTitle>Complainant&apos;s Evidence</CardTitle>
          <CardDescription></CardDescription>
        </CardHeader>
        <CardContent className="flex gap-2">
          {data.evidence.map((evi, i) => (
            <div key={i} className="rounded-lg bg-gray-950 p-4 text-center text-gray-50 w-40">
              <FileIcon className="mx-auto h-8 w-8" />
              <p className="mt-2 text-sm font-medium">{evi.label}</p>
            </div>
          ))}
        </CardContent>
      </Card>
      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Actions</CardTitle>
        </CardHeader>
        <CardContent asChild>
          <form onSubmit={onFilesSubmit}>
            <input type="hidden" name="dispute_id" value={data.id} />
            {files.map((file, i) => (
              <div key={i}>
                <span>{file.name}</span>
                <Button variant="destructive" onClick={() => removeFile(i)}>
                  Remove
                </Button>
              </div>
            ))}
            <Input type="file" placeholder="shadcn" multiple onChange={onFilesChange} />
            <Button disabled={files.length == 0} type="submit">
              Upload
            </Button>
          </form>
        </CardContent>
      </Card>
    </>
  );
}
