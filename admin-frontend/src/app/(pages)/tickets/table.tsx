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
import { createContext, useContext, useEffect, useState } from "react";
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

export interface TicketFilters {
  search?: string;
  page: number;
  filters: TicketFilter[];
}

const TicketContext = createContext<TicketFilters>({ filters: [], page: 0 });
export const TicketProvider = TicketContext.Provider;

export function TicketTable() {
  const filters = useContext(TicketContext);
  const { data, error, isPending } = useQuery({
    queryKey: ["ticketTable", filters],
    queryFn: () =>
      getTicketSummaries({
        search: filters.search,
        filter: filters.filters,
        limit: PAGE_SIZE,
        offset: PAGE_SIZE * filters.page,
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
}: {
  onValueChange?: (page: number) => void;
}) {
  const filters = useContext(TicketContext);
  const query = useQuery<TicketListResponse>({
    queryKey: ["ticketsTable", filters],
  });

  const [current, setCurrent] = useState(filters.page);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    setCurrent(filters.page);
  }, [filters.page]);

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
