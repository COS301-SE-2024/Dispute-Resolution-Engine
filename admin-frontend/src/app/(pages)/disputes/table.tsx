"use client";

import {
  TableCell,
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getDisputeList } from "@/lib/api/dispute";
import { useQuery } from "@tanstack/react-query";
import { AdminDisputesResponse, DisputeFilter, Filter, type AdminDispute } from "@/lib/types";

import { StatusBadge } from "@/components/admin/status-dropdown";
import { LinkIcon } from "lucide-react";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useToast } from "@/lib/hooks/use-toast";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
} from "@/components/ui/pagination";
import { PAGE_SIZE } from "@/lib/constants";

export function DisputeTable({ filters, page = 0 }: { filters?: DisputeFilter[]; page: number }) {
  const { data, error, isPending } = useQuery({
    queryKey: ["disputeTable", filters, page],
    queryFn: () =>
      getDisputeList({
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
          <TableHead className="">Title</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Workflow</TableHead>
          <TableHead className="w-[150px] text-center">Date Filed</TableHead>
          <TableHead className="w-[150px] text-center">Date Resolved</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {(data?.total ?? 0) > 0
          ? data?.disputes.map((dispute) => <DisputeRow key={dispute.id} {...dispute} />)
          : !isPending && (
              <TableRow>
                <TableCell>No results</TableCell>
              </TableRow>
            )}
      </TableBody>
    </Table>
  );
}

function DisputeRow(props: AdminDispute) {
  return (
    <TableRow className="text-nowrap truncate">
      <TableCell className="font-medium">
        <Link href={{ pathname: "/disputes", query: { id: props.id } }}>{props.title}</Link>
      </TableCell>
      <TableCell>
        <StatusBadge value={props.status} />
      </TableCell>
      <TableCell>
        <Link
          className="flex gap-1 items-center hover:underline text-nowrap"
          href={`/workflows/${props.workflow.id}`}
        >
          <span>{props.workflow.title}</span>
          <LinkIcon size="0.8rem" />
        </Link>
      </TableCell>
      <TableCell className="text-center">{props.date_filed}</TableCell>
      <TableCell className="text-center">{props.date_resolved ?? "-"}</TableCell>
    </TableRow>
  );
}

export function DisputePager({
  onValueChange = () => {},
  filters,
  page = 0,
}: {
  onValueChange?: (page: number) => void;
  filters?: DisputeFilter[];
  page?: number;
}) {
  const query = useQuery<AdminDisputesResponse>({
    queryKey: ["disputeTable", filters, page],
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
