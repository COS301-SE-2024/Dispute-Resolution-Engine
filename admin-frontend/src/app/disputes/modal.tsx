import { Button } from "@/components/ui/button";

import { DialogClose, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { X } from "lucide-react";
import Sidebar from "@/components/admin/sidebar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import Evidence from "./evidence";
import { Label } from "@/components/ui/label";
import { StatusBadge } from "@/components/admin/status-dropdown";

export default function DisputeDetails({ open = false }: { open?: boolean }) {
  return (
    <Sidebar open={open} className="p-6 md:pl-8 rounded-l-2xl flex flex-col">
      <DialogHeader className="grid grid-cols-[1fr_auto] gap-2 border-b pb-6 mb-6 border-primary-500/50 space-y-0 items-center">
        <DialogTitle className="p-2">
          Mr. Plastic CC vs. Mr. Plastic & Mining Promotional Goods
        </DialogTitle>
        <div className="flex justify-end items-start">
          <DialogClose asChild>
            <Button variant="ghost" className="rounded-full aspect-square p-2 m-0">
              <X />
            </Button>
          </DialogClose>
        </div>
        <div className="flex gap-2 items-center">
          <StatusBadge variant="waiting" dropdown>
            Awaiting respondent
          </StatusBadge>
          <span>Filed 2 days ago</span>
        </div>

        <p>Case Number: ZA2007-0001</p>
      </DialogHeader>
      <div className="overflow-y-auto grow space-y-6 pr-3">
        <Card>
          <CardHeader>
            <CardTitle>Overview</CardTitle>
            <CardDescription>Description here</CardDescription>
          </CardHeader>
          <CardContent>
            <h4 className="mb-1">Evidence</h4>
            <ul className="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-2">
              <Evidence id="0" label="evidence.pdf" url="https://google.com" date="today" />
              <Evidence id="1" label="evidence.pdf" url="https://google.com" date="today" />
              <Evidence id="2" label="evidence.pdf" url="https://google.com" date="today" />
              <Evidence id="3" label="evidence.pdf" url="https://google.com" date="today" />
              <Evidence id="4" label="evidence.pdf" url="https://google.com" date="today" />
              <Evidence id="5" label="evidence.pdf" url="https://google.com" date="today" />
            </ul>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="grid grid-cols-1 md:grid-cols-2 gap-3">
            <section className="space-y-3">
              <CardTitle>Complainant</CardTitle>
              <div>
                <Label>Name</Label>
                <Input disabled value="John doe" />
              </div>

              <div>
                <Label>Email</Label>
                <Input disabled value="email" />
              </div>

              <div>
                <Label>Address</Label>
                <Textarea disabled className="resize-none" value="Hey" />
              </div>
            </section>
            <section className="space-y-3">
              <CardTitle>Respondent</CardTitle>
              <div>
                <Label>Name</Label>
                <Input disabled value="John doe" />
              </div>

              <div>
                <Label>Email</Label>
                <Input value="email" />
              </div>

              <div>
                <Label>Address</Label>
                <Textarea disabled className="resize-none" value="Hey" />
              </div>
            </section>
          </CardContent>
        </Card>
      </div>
    </Sidebar>
  );
}
