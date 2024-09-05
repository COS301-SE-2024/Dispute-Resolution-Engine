"use client";

import { Button } from "@/components/ui/button";

import { DialogClose, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { X } from "lucide-react";
import Sidebar from "@/components/admin/sidebar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import Evidence from "./evidence";
import { Label } from "@/components/ui/label";
import { StatusBadge, StatusDropdown } from "@/components/admin/status-dropdown";
import { type UserDetails, type DisputeDetails, DisputeStatus } from "@/lib/types/dispute";
import { changeDisputeStatus } from "@/lib/api/dispute";

export default function DisputeDetails({ details }: { details?: DisputeDetails }) {
  async function changeStatus(status: DisputeStatus) {
    console.log(await changeDisputeStatus(details!.id, status));
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
            <StatusDropdown onSelect={changeStatus}>
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
              <ul className="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-2">
                {details.evidence.length > 0 ? (
                  details.evidence.map((evidence) => (
                    <Evidence
                      key={evidence.id}
                      id={evidence.id}
                      label={evidence.label}
                      url="https://google.com"
                      date={evidence.submitted_at}
                    />
                  ))
                ) : (
                  <li className="text-sm text-dark/50 dark:text-white/50">
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
        <Input value={d.email} />
      </div>

      <div>
        <Label>Address</Label>
        <Textarea disabled value={d.address} />
      </div>
    </>
  );
}
