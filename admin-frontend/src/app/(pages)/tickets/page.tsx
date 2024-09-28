"use client";

import PageHeader from "@/components/admin/page-header";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { TicketProvider, TicketsPager, TicketTable } from "./table";
import { z } from "zod";
import TicketDetails from "./details";
import { type TicketFilter } from "@/lib/types/tickets";
import SearchBar from "@/components/admin/search";
import { Button } from "@/components/ui/button";
import { FilterIcon } from "lucide-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import TicketFilters from "./ticket-filters";

const searchSchema = z.object({
  id: z.coerce.number().optional(),
});

export default function Tickets({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    throw new Error(JSON.stringify(searchError));
  }

  const [client] = useState(new QueryClient());
  const [filter, setFilter] = useState<TicketFilter[]>([]);
  const [page, setPage] = useState<number>(0);
  const [search, setSearch] = useState<string | undefined>();

  function changeSearch(search: string | undefined) {
    setSearch(search);
    setPage(0);
  }
  function changeFilter(filters: TicketFilter[]) {
    setFilter(filters);
    setPage(0);
  }

  return (
    <QueryClientProvider client={client}>
      {params.id && <TicketDetails ticketId={params.id} />}
      <TicketProvider value={{ filters: filter, page, search }}>
        <div className="flex flex-col">
          <PageHeader label="Tickets" />
          <div className="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
            <SearchBar placeholder="Search tickets..." onValueChange={changeSearch} timeout={500} />
            <TicketFilters onValueChange={changeFilter} value={filter}>
              <Button variant="ghost" className="gap-2">
                <FilterIcon />
                <span>Filter by</span>
              </Button>
            </TicketFilters>
          </div>
          <main className="overflow-auto p-5 grow">
            <Card>
              <CardContent>
                <TicketTable />
              </CardContent>
              <CardFooter>
                <TicketsPager onValueChange={setPage} />
              </CardFooter>
            </Card>
          </main>
        </div>
      </TicketProvider>
    </QueryClientProvider>
  );
}
