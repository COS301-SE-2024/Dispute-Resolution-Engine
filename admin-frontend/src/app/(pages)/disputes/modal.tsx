"use client";

import { Button } from "@/components/ui/button";

import { DialogClose, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { X } from "lucide-react";
import Sidebar from "@/components/admin/sidebar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { StatusBadge, StatusDropdown } from "@/components/admin/status-dropdown";
import { type UserDetails, type DisputeDetails, DisputeStatus } from "@/lib/types/dispute";
import { changeDisputeStatus, deleteEvidence } from "@/lib/api/dispute";
import { useToast } from "@/lib/hooks/use-toast";
import { useState } from "react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Download, EllipsisVertical, FileText, Trash } from "lucide-react";
import Link from "next/link";

export default function DisputeDetails({ details }: { details?: DisputeDetails }) {
  const { toast } = useToast();

  const [statusDisabled, setStatusDisabled] = useState(false);

  async function changeStatus(status: DisputeStatus) {
    setStatusDisabled(true);
    const { data, error } = await changeDisputeStatus(details!.id, status);
    setStatusDisabled(false);
    if (data) {
      toast({
        title: "Status updated successfully",
      });
    } else if (error) {
      toast({
        variant: "error",
        title: "Something went wrong",
        description: error,
      });
    }
  }

  async function deleteEvi(id: string) {
    const { data, error } = await deleteEvidence(details!.id, id);
    if (data) {
      toast({
        title: "Evidence removed",
      });
    } else if (error) {
      toast({
        variant: "error",
        title: "Something went wrong",
        description: error,
      });
    }
  }

  return (
    details && (
      <Sidebar open className="p-6 md:pl-8 rounded-l-2xl flex flex-col">
        <DialogHeader className="grid grid-cols-[1fr_auto] gap-2 border-b pb-6 mb-6 border-primary-500/50 space-y-0 items-center">
          <DialogTitle className="p-2">{details.title}</DialogTitle>
          <div className="flex justify-end items-start">
            <DialogClose asChild>
              <Button variant="ghost" className="rounded-full aspect-square p-2 m-0">
                <X />
              </Button>
            </DialogClose>
          </div>
          <div className="flex gap-2 items-center">
            <StatusDropdown onSelect={changeStatus} disabled={statusDisabled}>
              <StatusBadge variant="waiting" dropdown>
                {details.status}
              </StatusBadge>
            </StatusDropdown>
            <span>{details.date_filed}</span>
          </div>

          <p>Case Number: {details.id}</p>
        </DialogHeader>
        <div className="overflow-y-auto grow space-y-6 pr-3">
          <Card>
            <CardHeader>
              <CardTitle>Overview</CardTitle>
              <CardDescription>{details.description}</CardDescription>
            </CardHeader>
            <CardContent>
              <h4 className="mb-1">Evidence</h4>
              <ul className="grid grid-cols-[repeat(auto-fill,minmax(300px,1fr))] gap-2">
                {details.evidence.length > 0 ? (
                  details.evidence.map((evidence) => (
                    <Evidence
                      onDelete={deleteEvi}
                      key={evidence.id}
                      id={evidence.id}
                      label={evidence.label}
                      url="https://google.com"
                      date={evidence.submitted_at}
                    />
                  ))
                ) : (
                  <li className="text-sm text-black/50 dark:text-white/50">
                    No evidence has been submitted
                  </li>
                )}
              </ul>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <section className="space-y-3">
                <CardTitle>Complainant</CardTitle>
                <UserDetails {...details.complainant} />
              </section>
              <section className="space-y-3">
                <CardTitle>Respondent</CardTitle>
                <UserDetails {...details.respondent} />
              </section>
            </CardContent>
          </Card>
        </div>
      </Sidebar>
    )
  );
}
function UserDetails(d: UserDetails) {
  return (
    <>
      <div>
        <Label>Name</Label>
        <Input disabled value={d.name} />
      </div>

      <div>
        <Label>Email</Label>
        <Input disabled value={d.email} />
      </div>

      <div>
        <Label>Address</Label>
        <Textarea disabled value={d.address} />
      </div>
    </>
  );
}

function Evidence({
  id,
  label,
  url,
  date,
  onDelete,
}: {
  id: string;
  label: string;
  url: string;
  date: string;
  onDelete: (id: string) => void;
}) {
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
          <DropdownMenuItem asChild>
            <Link href={url}>
              <Download className="mr-2" />
              <span>Download file</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem className="text-red-500" onSelect={() => onDelete(id)}>
            <Trash className="mr-2" />
            <span>Delete evidence</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </li>
  );
}
