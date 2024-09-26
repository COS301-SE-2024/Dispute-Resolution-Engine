"use client";

import {
  TableCell,
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { keepPreviousData, useQuery } from "@tanstack/react-query";

import Link from "next/link";
import { createContext, useContext, useEffect, useState } from "react";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
} from "@/components/ui/pagination";
import { cn } from "@/lib/utils";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { PAGE_SIZE, WORKFLOW_LIST_KEY } from "@/lib/constants";
import { getWorkflowList } from "@/lib/api/workflow";
import { WorkflowListResponse, WorkflowSummary } from "@/lib/types/workflow";

export interface WorkflowFilters {
  search?: string;
  page: number;
}

const WorkflowContext = createContext<WorkflowFilters>({ page: 0 });
export const WorkflowProvider = WorkflowContext.Provider;

export function WorkflowTable() {
  const filters = useContext(WorkflowContext);
  const { data, error, isPending } = useQuery({
    queryKey: [WORKFLOW_LIST_KEY, filters],
    queryFn: () =>
      getWorkflowList({
        search: filters.search,
        limit: PAGE_SIZE,
        offset: PAGE_SIZE * filters.page,
      }),
    placeholderData: keepPreviousData,
  });

  useErrorToast(error, "Failed to fetch ticket list");

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead className="">Created by</TableHead>
          <TableHead className="w-[150px] text-center">Date created</TableHead>
          <TableHead className="w-[150px] text-center">Last updated</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {(data?.total ?? 0) > 0
          ? data?.workflows.map((t) => <WorkflowRow key={t.id} {...t} />)
          : !isPending && (
              <TableRow>
                <TableCell>No results</TableCell>
              </TableRow>
            )}
      </TableBody>
    </Table>
  );
}

function WorkflowRow(props: WorkflowSummary) {
  return (
    <TableRow className="text-nowrap truncate">
      <TableCell className="font-medium hover:underline">
        <Link href={`/workflows/designer/${props.id}`}>{props.name}</Link>
      </TableCell>
      <TableCell>{props.author.full_name}</TableCell>
      <TableCell className="text-center">{props.date_created}</TableCell>
      <TableCell className="text-center">{props.last_updated}</TableCell>
    </TableRow>
  );
}

export function WorkflowsPager({
  onValueChange = () => {},
}: {
  onValueChange?: (page: number) => void;
}) {
  const filters = useContext(WorkflowContext);
  const query = useQuery<WorkflowListResponse>({
    queryKey: [WORKFLOW_LIST_KEY, filters],
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
    <Pagination className={cn("w-full", query.isPending && "opacity-50")}>
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
  );
}
