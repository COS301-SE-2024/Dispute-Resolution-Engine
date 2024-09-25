"use client";

import PageHeader from "@/components/admin/page-header";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { TicketProvider, TicketsPager, TicketTable } from "./table";
import { z } from "zod";
import TicketDetails from "./details";
import { Ticket, TicketFilter } from "@/lib/types/tickets";
import SearchBar from "@/components/admin/search";
import { Button } from "@/components/ui/button";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
} from "@/components/ui/pagination";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { FilterIcon } from "lucide-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { createContext } from "vm";

const searchSchema = z.object({
  id: z.string().optional(),
});

const ticket = {
  id: "0",
  user: { id: "0", full_name: "John Doe" },
  date_created: "2023-09-01T10:30:00Z",
  subject: "Unable to access account",
  status: "Open",
  body: "I cannot log in to my account. It keeps saying incorrect password.",
  messages: [
    {
      id: "0",
      user: { id: "0", full_name: "John Doe" },
      date_sent: "2023-09-01T10:31:00Z",
      message: "I cannot log in to my account. It keeps saying incorrect password.",
    },
  ],
} satisfies Ticket;

export default function Tickets({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    throw new Error(JSON.stringify(searchError));
  }

  const client = new QueryClient();

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
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="ghost" className="gap-2">
                  <FilterIcon />
                  <span>Filter by</span>
                </Button>
              </PopoverTrigger>
              <PopoverContent className="grid gap-x-2 gap-y-3 grid-cols-[auto_1fr] items-center">
                <strong className="col-span-2">Filter</strong>
              </PopoverContent>
            </Popover>
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
