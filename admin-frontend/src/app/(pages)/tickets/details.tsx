"use client";

import Sidebar from "@/components/admin/sidebar";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { TicketStatusBadge } from "@/components/admin/status-badge";
import { Ticket, TicketMessage, TicketStatus } from "@/lib/types/tickets";
import SidebarHeader from "@/components/sidebar/header";
import { TicketStatusDropdown } from "@/components/admin/status-dropdown";
import { useToast } from "@/lib/hooks/use-toast";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { addTicketMessage, changeTicketStatus, getTicketDetails } from "@/lib/api/tickets";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { FormEvent } from "react";
import { TICKET_DETAILS_KEY, TICKET_LIST_KEY } from "@/lib/constants";
import Link from "next/link";

export default function TicketDetails({ ticketId }: { ticketId: number }) {
  const { data, error } = useQuery({
    queryKey: [TICKET_DETAILS_KEY, ticketId],
    queryFn: async () => getTicketDetails(ticketId),
  });
  useErrorToast(error, "Failed to fetch ticket details");

  const client = useQueryClient();
  const { toast } = useToast();
  const status = useMutation({
    mutationFn: (status: TicketStatus) => changeTicketStatus(ticketId, status),
    onSuccess: (data, variables) => {
      client.refetchQueries({ queryKey: [TICKET_LIST_KEY] });
      client.setQueryData([TICKET_DETAILS_KEY, ticketId], (old: Ticket) => ({
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

  const message = useMutation({
    mutationFn: (message: string) => addTicketMessage(ticketId, message),
    onSuccess: async (data, variables) => {
      client.invalidateQueries({ queryKey: [TICKET_LIST_KEY] });
      client.setQueryData([TICKET_DETAILS_KEY, ticketId], (old: Ticket) => ({
        ...old,
        messages: [...old.messages, data],
      }));
    },
    onError: (error) => {
      toast({
        variant: "error",
        title: "Failed to send message",
        description: error?.message,
      });
    },
  });

  function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const formdata = new FormData(e.currentTarget!);
    const data = formdata.get("message")! as string;
    message.mutate(data);
  }

  return (
    <Sidebar open className="p-6 md:pl-8 rounded-l-2xl flex flex-col">
      {data && (
        <>
          <SidebarHeader title={data.subject} className="flex gap-2 items-center flex-wrap">
            <TicketStatusDropdown onSelect={status.mutate}>
              <TicketStatusBadge variant={data.status} dropdown>
                {data.status}
              </TicketStatusBadge>
            </TicketStatusDropdown>
            <span className="grow">{data.date_created}</span>
            <Link href={{ pathname: "/disputes", query: { id: data.dispute_id } }}>
              Go to dispute
            </Link>
          </SidebarHeader>
          <div className="overflow-y-auto grow space-y-6 pr-3">
            <Card>
              <CardContent>{data.body}</CardContent>
            </Card>
            {data.messages.map((msg) => (
              <TicketMessageCard key={msg.id} {...msg} />
            ))}
            <Card asChild>
              <form onSubmit={onSubmit}>
                <CardHeader>
                  <CardTitle>Send a message</CardTitle>
                  <CardDescription>Enter a message to reply to the ticket</CardDescription>
                </CardHeader>
                <CardContent>
                  <Textarea name="message" />
                </CardContent>
                <CardFooter>
                  <Button>Send</Button>
                </CardFooter>
              </form>
            </Card>
          </div>
        </>
      )}
    </Sidebar>
  );
}

function TicketMessageCard(msg: TicketMessage) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{msg.user.full_name}</CardTitle>
        <CardDescription>Sent at {msg.date_sent}</CardDescription>
      </CardHeader>
      <CardContent>{msg.message}</CardContent>
    </Card>
  );
}
