import { Filter, Search } from "lucide-react";
import { z } from "zod";

import PageHeader from "@/components/admin/page-header";
import StatusFilter from "@/components/dispute/status-filter";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@/components/ui/table";

import { getDisputeDetails, getDisputeList } from "@/lib/api/dispute";
import { type DisputeDetails } from "@/lib/types/dispute";

import Details from "./modal";
import DisputeRow from "./row";

const searchSchema = z.object({
  id: z.string().optional(),
});

export default async function Disputes({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    throw new Error(JSON.stringify(searchError));
  }

  const { data, error } = await getDisputeList({});
  if (error) {
    throw new Error(error);
  }

  let details: DisputeDetails | undefined = undefined;
  if (params.id) {
    const { data, error } = await getDisputeDetails(params.id);
    if (error) {
      throw new Error(error);
    }
    details = data;
  }

  return (
    <>
      <Details details={details} />
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
              <Search size={20} />
            </div>
          </div>
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="ghost" className="gap-2">
                <Filter />
                <span>Filter by</span>
              </Button>
            </PopoverTrigger>
            <PopoverContent className="grid gap-x-2 gap-y-3 grid-cols-[auto_1fr] items-center">
              <strong className="col-span-2">Filter</strong>

              <label>Status</label>
              <StatusFilter />
              <label>Workflow</label>
              <StatusSelect />
              <div className="col-span-2 flex flex-end">
                <Button className="ml-auto">Apply</Button>
              </div>
            </PopoverContent>
          </Popover>
        </div>
        <main className="overflow-auto p-5 grow">
          <Card>
            <CardContent>
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
                  {data!.map((dispute) => (
                    <DisputeRow key={dispute.id} {...dispute} />
                  ))}
                </TableBody>
              </Table>
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

function StatusSelect() {
  return (
    <Select>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select a fruit" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Fruits</SelectLabel>
          <SelectItem value="apple">Apple</SelectItem>
          <SelectItem value="banana">Banana</SelectItem>
          <SelectItem value="blueberry">Blueberry</SelectItem>
          <SelectItem value="grapes">Grapes</SelectItem>
          <SelectItem value="pineapple">Pineapple</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
