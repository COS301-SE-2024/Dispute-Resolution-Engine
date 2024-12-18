"use client";

import PageHeader from "@/components/admin/page-header";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { TicketProvider } from "../tickets/table";
import { WorkflowProvider, WorkflowsPager, WorkflowTable } from "./table";
import SearchBar from "@/components/admin/search";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { PlusIcon } from "lucide-react";
import Link from "next/link";

export default function Workflows() {
  const [client] = useState(new QueryClient());
  const [page, setPage] = useState<number>(0);
  const [search, setSearch] = useState<string | undefined>();

  function changeSearch(search: string | undefined) {
    setSearch(search);
    setPage(0);
  }

  return (
    <QueryClientProvider client={client}>
      <WorkflowProvider value={{ page, search }}>
        <div className="flex flex-col">
          <PageHeader label="Workflows" />
          <div className="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
            <SearchBar placeholder="Search tickets..." onValueChange={changeSearch} timeout={500} />
            <Button className="gap-2" asChild>
              <Link href="/workflows/designer">
                <PlusIcon size="1rem" />
                <span>New workflow</span>
              </Link>
            </Button>
          </div>
          <main className="overflow-auto p-5 grow">
            <Card>
              <CardContent>
                <WorkflowTable />
              </CardContent>
              <CardFooter>
                <WorkflowsPager onValueChange={setPage} />
              </CardFooter>
            </Card>
          </main>
        </div>
      </WorkflowProvider>
    </QueryClientProvider>
  );
}
