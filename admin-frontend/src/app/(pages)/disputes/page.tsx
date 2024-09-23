"use client";
import { FilterIcon, SearchIcon } from "lucide-react";
import { z } from "zod";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import PageHeader from "@/components/admin/page-header";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

import { type Filter, type DisputeDetails } from "@/lib/types";

import Details from "./modal";
import DisputeFilter from "./dispute-filter";
import DisputeTable from "./table";
import { useState } from "react";

const searchSchema = z.object({
  id: z.string().optional(),
});

export default function Disputes({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    throw new Error(JSON.stringify(searchError));
  }

  const client = new QueryClient();

  const [filter, setFilter] = useState<Filter[]>([]);
  return (
    <QueryClientProvider client={client}>
      {params.id && <Details id={params.id} />}
      <div className="flex flex-col">
        <PageHeader label="Disputes" />
        <div className="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
          <div className="grid grid-cols-[auto_1fr] items-center grow">
            <input
              type="search"
              className="col-span-2 p-5 bg-transparent  col-start-1 row-start-1 pl-12"
              placeholder="Search disputes..."
            />
            <div className="p-5 row-start-1 col-start-1 pointer-events-none">
              <SearchIcon size={20} />
            </div>
          </div>

          <DisputeFilter onValueChange={setFilter}>
            <Button variant="ghost" className="gap-2">
              <FilterIcon />
              <span>Filter by</span>
            </Button>
          </DisputeFilter>
        </div>
        <main className="overflow-auto p-5 grow">
          <Card>
            <CardContent>
              <DisputeTable filters={filter} />
            </CardContent>
            <CardFooter>
              <Pagination className="w-full">
                <PaginationContent className="w-full">
                  <PaginationItem>
                    <PaginationPrevious href="#" />
                  </PaginationItem>
                  <div className="grow" />
                  <PaginationItem>
                    <PaginationNext href="#" />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </CardFooter>
          </Card>
        </main>
      </div>
    </QueryClientProvider>
  );
}
