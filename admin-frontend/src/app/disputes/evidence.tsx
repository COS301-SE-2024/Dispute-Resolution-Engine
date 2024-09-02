"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Download, EllipsisVertical, FileText, Trash } from "lucide-react";
import Link from "next/link";

export default function Evidence({
  id,
  label,
  url,
  date,
}: {
  id: string;
  label: string;
  url: string;
  date: string;
}) {
  function downloadEvidence() {
    alert("Download evidence");
    // TODO: Integrate with API endpoint using `id` and `url` props
  }
  function deleteEvidence() {
    alert("Delete evidence");
    // TODO: Integrate with API endpoint using `id` prop
  }

  return (
    <li className="grid grid-cols-[auto_1fr_auto] gap-2 items-center px-3 py-2 border border-primary-500/30 rounded-md">
      <FileText className="stroke-primary-500" size="1.7rem" />
      <div>
        <span className="truncate">{label}</span> <br />
        <span className="truncate opacity-50">{date}</span>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="rounded-full p-2">
            <EllipsisVertical />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onSelect={() => downloadEvidence()}>
            <Download className="mr-2" />
            <span>Download file</span>
          </DropdownMenuItem>
          <DropdownMenuItem className="text-red-500" onSelect={() => deleteEvidence()}>
            <Trash className="mr-2" />
            <span>Delete evidence</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </li>
  );
}
