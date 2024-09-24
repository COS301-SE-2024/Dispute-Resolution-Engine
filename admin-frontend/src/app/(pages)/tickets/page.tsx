import PageHeader from "@/components/admin/page-header";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { TicketTable } from "./table";
import { z } from "zod";
import TicketDetails from "./details";
import { Ticket } from "@/lib/types/tickets";
import SearchBar from "@/components/admin/search";
import { Button } from "@/components/ui/button";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
} from "@/components/ui/pagination";
import { TableHeader, TableRow, TableHead, TableBody } from "@/components/ui/table";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { FilterIcon, Table } from "lucide-react";
import DisputeRow from "../disputes/row";

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

  return (
    <>
      {params.id && <TicketDetails details={ticket} />}
      <div className="flex flex-col">
        <PageHeader label="Tickets" />
        <div className="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
          <SearchBar placeholder="Search tickets..." />
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
    </>
  );
}
