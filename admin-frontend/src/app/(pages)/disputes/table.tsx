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
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Filter, type AdminDispute } from "@/lib/types";

import { StatusBadge } from "@/components/admin/status-dropdown";
import { LinkIcon } from "lucide-react";
import Link from "next/link";
import { unwrapResult } from "@/lib/utils";

export default function DisputeTable({ filters }: { filters?: Filter[] }) {
  const { data, error } = useQuery({
    queryKey: ["disputeTable", filters],
    queryFn: async () =>
      unwrapResult(
        getDisputeList({
          filter: filters,
        })
      ),
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
        {data?.map((dispute) => (
          <DisputeRow key={dispute.id} {...dispute} />
        ))}
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
        <StatusBadge variant="active">{props.status}</StatusBadge>
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
