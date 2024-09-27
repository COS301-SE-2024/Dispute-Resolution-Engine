"use client";
import { FilterIcon, SearchIcon } from "lucide-react";
import { z } from "zod";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import PageHeader from "@/components/admin/page-header";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";

import { DisputeFilter } from "@/lib/types";

import Details from "./modal";
import SearchBar from "@/components/admin/search";
import { DisputePager, DisputeTable } from "./table";
import { useState } from "react";
import DisputeFilters from "./dispute-filter";

const searchSchema = z.object({
  id: z.string().optional(),
});

export default function Disputes({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    throw new Error(JSON.stringify(searchError));
  }

  const [client] = useState(new QueryClient());

  const [filter, setFilter] = useState<DisputeFilter[]>([]);
  const [page, setPage] = useState<number>(0);

  return (
    <QueryClientProvider client={client}>
      {params.id && <Details id={params.id} />}
      <div className="flex flex-col">
        <PageHeader label="Disputes" />
        <div className="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
          <SearchBar placeholder="Search disputes..." />

          <DisputeFilters
            value={filter}
            onValueChange={(f) => {
              setFilter(f);
              setPage(0);
            }}
          >
            <Button variant="ghost" className="gap-2">
              <FilterIcon />
              <span>Filter by</span>
            </Button>
          </DisputeFilters>
        </div>
        <main className="overflow-auto p-5 grow">
          <Card>
            <CardContent>
              <DisputeTable page={page} filters={filter} />
            </CardContent>
            <CardFooter>
              <DisputePager page={page} filters={filter} onValueChange={setPage} />
            </CardFooter>
          </Card>
        </main>
      </div>
    </QueryClientProvider>
  );
}
