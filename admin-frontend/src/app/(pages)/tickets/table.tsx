"use client";

import {
  TableCell,
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useQuery } from "@tanstack/react-query";
import { AdminDisputesResponse, Filter } from "@/lib/types";

import { TicketStatusBadge } from "@/components/admin/status-badge";
import { LinkIcon } from "lucide-react";
import Link from "next/link";
import { createContext, useEffect, useState } from "react";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
} from "@/components/ui/pagination";
import { Ticket, TicketFilter, TicketListResponse, TicketSummary } from "@/lib/types/tickets";
import { getTicketSummaries } from "@/lib/api/tickets";
import { useToast } from "@/lib/hooks/use-toast";

const PAGE_SIZE = 10;

export function TicketTable({
  filters,
  page = 0,
  search,
}: {
  filters: TicketFilter[];
  page?: number;
  search?: string;
}) {
  const { data, error, isPending } = useQuery({
    queryKey: ["ticketTable", filters, page, search],
    queryFn: () =>
      getTicketSummaries({
        search: search,
        filter: filters,
        limit: PAGE_SIZE,
        offset: PAGE_SIZE * page,
      }),
  });

  const { toast } = useToast();
  useEffect(() => {
    if (error) {
      toast({
        variant: "error",
        title: "Failed to fetch dispute list",
        description: error?.message,
      });
    }
  });
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Subject</TableHead>
          <TableHead className="">Created by</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="w-[150px] text-center">Date Created</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {(data?.total ?? 0) > 0
          ? data?.tickets.map((t) => <TicketRow key={t.id} {...t} />)
          : !isPending && (
              <TableRow>
                <TableCell>No results</TableCell>
              </TableRow>
            )}
      </TableBody>
    </Table>
  );
}

function TicketRow(props: TicketSummary) {
  return (
    <TableRow className="text-nowrap truncate">
      <TableCell className="font-medium">
        <Link href={{ pathname: "/tickets", query: { id: props.id } }}>{props.subject}</Link>
      </TableCell>
      <TableCell className="font-medium">{props.user.full_name}</TableCell>
      <TableCell>
        <TicketStatusBadge variant={props.status}>{props.status}</TicketStatusBadge>
      </TableCell>
      <TableCell className="text-center">{props.date_created}</TableCell>
    </TableRow>
  );
}

export function TicketsPager({
  onValueChange = () => {},
  filters,
  page = 0,
  search,
}: {
  onValueChange?: (page: number) => void;
  filters?: TicketFilter[];
  page?: number;
  search?: string;
}) {
  const query = useQuery<TicketListResponse>({
    queryKey: ["ticketsTable", filters, page, search],
  });

  const [current, setCurrent] = useState(page);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    setCurrent(page);
  }, [page]);
  useEffect(() => {
    if (!query.data) {
      setTotal(0);
    } else {
      setTotal(Math.ceil(query.data.total / PAGE_SIZE));
    }
  }, [query.data]);

  function navigate(page: number) {
    setCurrent(page);
    onValueChange(page);
  }

  return (
    query.isSuccess && (
      <Pagination className="w-full">
        <PaginationContent className="w-full">
          <PaginationItem>
            <PaginationPrevious disabled={current == 0} onClick={() => navigate(current - 1)} />
          </PaginationItem>
          <div className="grow flex justify-center items-center">
            {current + 1}/{total}
          </div>
          <PaginationItem>
            <PaginationNext disabled={current >= total - 1} onClick={() => navigate(current + 1)} />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    )
  );
}
