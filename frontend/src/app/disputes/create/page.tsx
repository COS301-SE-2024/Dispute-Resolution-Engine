"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import Link from "next/link";
import { API_URL } from "@/lib/utils";

export default function CreateDispute() {
  return (
    <div className="grow overflow-y-auto flex flex-col">
      <header className="px-3 py-6 flex">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
          Create a Dispute
        </h1>
      </header>
      <Separator />
      <main className="grow overflow-y-auto p-5 space-y-2">
        <div>
          <Label>Title</Label>
          <Input placeholder="Title" />
        </div>
        <div className="space-y-2">
          <Label>Respondant Details</Label>
          <Input type="email" placeholder="Respondant's Email" />
          <Input type="tel" placeholder="Respondant's Telephone" />
        </div>
        <div>
          <Label>Summary</Label>
          <Textarea placeholder="Write a short description of the dispute..." />
        </div>
        <div>
          <Label>Evidence</Label>
          <Input placeholder="Evidence" type="file" />
        </div>
      </main>
      <footer className="p-2 flex justify-between">
        <Button asChild onClick={createDisputeAPI}>
          <Link href="/disputes">Create</Link>
        </Button>
        <Button variant={"destructive"} asChild>
          <Link href="/disputes">Cancel</Link>
        </Button>
      </footer>
    </div>
  );
}
async function createDisputeAPI(data: any) {
  const returnData = await fetch(`${API_URL}/createDispute`, {
    method: "POST",
    body: JSON.stringify({
      title: data.title,
      respondent: data.respondent,
      summary: data.summary,
      evidence: data.evidence,
    }),
  })
}