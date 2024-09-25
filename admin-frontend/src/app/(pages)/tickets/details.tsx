"use client";

import Sidebar from "@/components/admin/sidebar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { TicketStatusBadge } from "@/components/admin/status-badge";
import { Ticket, TicketStatus } from "@/lib/types/tickets";
import SidebarHeader from "@/components/sidebar/header";
import { TicketStatusDropdown } from "@/components/admin/status-dropdown";
import { useToast } from "@/lib/hooks/use-toast";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { changeTicketStatus, getTicketDetails } from "@/lib/api/tickets";
import { useErrorToast } from "@/lib/hooks/use-query-toast";

const DETAILS_KEY = "ticketDetails";

export default function TicketDetails({ ticketId }: { ticketId: string }) {
  const { data, error } = useQuery({
    queryKey: [DETAILS_KEY, ticketId],
    queryFn: async () => getTicketDetails(ticketId),
  });
  useErrorToast(error, "Failed to fetch ticket details");

  const client = useQueryClient();
  const { toast } = useToast();
  const status = useMutation({
    mutationFn: (status: TicketStatus) => changeTicketStatus(ticketId, status),
    onSuccess: (data, variables) => {
      client.setQueryData([DETAILS_KEY, ticketId], (old: Ticket) => ({
        ...old,
        status: variables,
      }));
      toast({
        title: "Status updated successfully",
      });
    },
    onError: (error) => {
      toast({
        variant: "error",
        title: "Status update failed",
        description: error?.message,
      });
    },
  });

  return (
    <Sidebar open className="p-6 md:pl-8 rounded-l-2xl flex flex-col">
      {data && (
        <>
          <SidebarHeader title={data.subject} className="flex gap-2 items-center">
            <TicketStatusDropdown onSelect={status.mutate}>
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
