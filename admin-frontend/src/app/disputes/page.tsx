import { Search, Filter, ChevronLeft, ChevronRight } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Card, CardContent, CardFooter } from "@/components/ui/card";

export default function Disputes() {
  return (
    <div className="flex flex-col">
      <header>
        <h2 className="p-5 font-bold tracking-wide text-xl border-b dark:border-primary-500/30 border-primary-500/20">
          Disputes
        </h2>
        <div class="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
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
          <Button variant="ghost" className="gap-2">
            <Filter />
            <span>Filter by</span>
          </Button>
        </div>
      </header>
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
                <TableRow>
                  <TableCell className="font-medium">INV001</TableCell>
                  <TableCell>Paid</TableCell>
                  <TableCell>Credit Card</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium">INV001</TableCell>
                  <TableCell>Paid</TableCell>
                  <TableCell>Credit Card</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium">INV001</TableCell>
                  <TableCell>Paid</TableCell>
                  <TableCell>Credit Card</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell className="font-medium">INV001</TableCell>
                  <TableCell>Paid</TableCell>
                  <TableCell>Credit Card</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                  <TableCell className="text-center">$250.00</TableCell>
                </TableRow>
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
  );
}
