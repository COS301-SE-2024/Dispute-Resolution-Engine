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
import { TicketStatusDropdown } from "@/components/admin/status-dropdown";
import { useToast } from "@/lib/hooks/use-toast";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getTicketDetails } from "@/lib/api/tickets";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { useEffect } from "react";

export default function TicketDetails({ ticketId }: { ticketId: string }) {
  // const { toast } = useToast();
  // const client = useQueryClient();
  const { data, error } = useQuery({
    queryKey: ["ticket"],
    queryFn: async () => getTicketDetails(ticketId),
  });
  useErrorToast(error, "Failed to fetch ticket details");

  // const status = useMutation({
  //   mutationFn: async (status: DisputeStatus) => {
  //     await unwrapResult(changeDisputeStatus(disputeId, status));
  //   },
  //   onSuccess: (data, variables) => {
  //     client.setQueryData(["dispute"], (old: DisputeDetailsResponse) => ({
  //       ...old,
  //       status: variables,
  //     }));
  //     toast({
  //       title: "Status updated successfully",
  //     });
  //   },
  //   onError: (error) => {
  //     toast({
  //       variant: "error",
  //       title: "Something went wrong",
  //       description: error?.message,
  //     });
  //   },
  // });

  return (
    <Sidebar open className="p-6 md:pl-8 rounded-l-2xl flex flex-col">
      {data && (
        <>
          <SidebarHeader title={data.subject} className="flex gap-2 items-center">
            <TicketStatusDropdown>
              <TicketStatusBadge variant={data.status} dropdown>
                {data.status}
              </TicketStatusBadge>
            </TicketStatusDropdown>
            <span>{data.date_created}</span>
          </SidebarHeader>
          <div className="overflow-y-auto grow space-y-6 pr-3">
            <Card>
              <CardContent>{data.body}</CardContent>
            </Card>
            {data.messages.map((msg) => (
              <Card key={msg.id}>
                <CardHeader>
                  <CardTitle>{msg.user.full_name}</CardTitle>
                  <CardDescription>Sent at {msg.date_sent}</CardDescription>
                </CardHeader>
                <CardContent>{msg.message}</CardContent>
              </Card>
            ))}
          </div>
        </>
      )}
    </Sidebar>
  );
}
