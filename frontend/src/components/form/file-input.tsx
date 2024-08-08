"use client";

import { ChangeEvent, HTMLAttributes, useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Trash } from "lucide-react";

export type FileInputProps = Pick<
  HTMLAttributes<HTMLInputElement>,
  Exclude<keyof HTMLAttributes<HTMLInputElement>, "onChange" | "multiple" | "type" | "name">
> & {
  onValueChange?: (files: File[]) => void;
};

export default function FileInput({ onValueChange = () => {}, ...props }: FileInputProps) {
  const [files, setFiles] = useState<File[]>([]);

  const onFilesChange = async (ev: ChangeEvent<HTMLInputElement>) => {
    const newFiles = [...files, ...ev.target.files!];
    setFiles(newFiles);
    onValueChange(newFiles);
    ev.target.value = "";
  };

  const removeFile = async (i: number) => {
    const newFiles = files.filter((_f, j) => i !== j);
    setFiles(newFiles);
    onValueChange(newFiles);
  };

  return (
    <ul className="space-y-3">
      {files.map((file, i) => (
        <li key={i} className="flex items-center gap-2">
          {/*
          Buttons trigger forms by default; the type="button" attribute is added to override that
          Source: https://stackoverflow.com/questions/932653/how-to-prevent-buttons-from-submitting-forms
          */}
          <Button
            type="button"
            variant="ghost"
            onClick={() => removeFile(i)}
            className="aspect-square p-1 justify-center rounded-full"
            title="Remove"
          >
            <Trash size="1rem" />
          </Button>
          <span>{file.name}</span>
        </li>
      ))}
      <li>
        <Input type="file" multiple onChange={onFilesChange} {...props} />
      </li>
    </ul>
  );
}
