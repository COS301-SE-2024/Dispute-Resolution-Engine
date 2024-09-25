"use client";

import { Button } from "@/components/ui/button";

import { DialogClose, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { X } from "lucide-react";
import Sidebar from "@/components/admin/sidebar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { TicketStatusBadge } from "@/components/admin/status-badge";
import { Ticket } from "@/lib/types/tickets";
import SidebarHeader from "@/components/sidebar/header";

export default function TicketDetails({ details }: { details: Ticket }) {
  return (
    <Sidebar open className="p-6 md:pl-8 rounded-l-2xl flex flex-col">
      <SidebarHeader title={details.subject} className="flex gap-2 items-center">
        <TicketStatusBadge variant={details.status} dropdown>
          {details.status}
        </TicketStatusBadge>
        <span>{details.date_created}</span>
      </SidebarHeader>
      <div className="overflow-y-auto grow space-y-6 pr-3">
        <Card>
          <CardContent>{details.body}</CardContent>
        </Card>
        {details.messages.map((msg) => (
          <Card key={msg.id}>
            <CardHeader>
              <CardTitle>{msg.user.full_name}</CardTitle>
              <CardDescription>Sent at {msg.date_sent}</CardDescription>
            </CardHeader>
            <CardContent>{msg.message}</CardContent>
          </Card>
        ))}
      </div>
    </Sidebar>
  );
}
