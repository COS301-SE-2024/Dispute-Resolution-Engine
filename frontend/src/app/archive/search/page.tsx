import { ExternalLink } from "lucide-react";
import Link from "next/link";
import { redirect } from "next/navigation";
import { z } from "zod";

import { searchArchive } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

const searchSchema = z.object({
  q: z.string().optional(),
  offset: z.coerce.number().default(0),
});

type SearchParams = z.infer<typeof searchSchema>;

function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <Card className="border-none">
      <CardHeader>
        <Link href="/archive/id">
          <CardTitle>{props.title}</CardTitle>
        </Link>
      </CardHeader>
      <CardContent asChild className="dark:text-white/50">
        <p>{props.summary}</p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <p>{props.resolution}</p>
        <Button>
          <ExternalLink size="1rem" className="mr-2" />
          Read More
        </Button>
      </CardFooter>
    </Card>
  );
}

function ErrorPage({ msg }: { msg: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-full gap-5 w-2/3 mx-auto">
      <h1 className="text-4xl font-bold tracking-wide">Oops, something went wrong :(</h1>
      <p className="dark:text-white/50">{msg}</p>
    </div>
  );
}

function pager(params: SearchParams, offset: number) {
  return { pathname: "/archive/search", query: { ...params, offset } };
}

export default async function ArchiveSearch({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    return <ErrorPage msg={JSON.stringify(searchError)} />;
  }

  if (!params.q) {
    redirect("/archive");
  }

  const { data, error } = await searchArchive({
    search: params.q,
  });

  // TODO: replace this with response information
  const total = 10;

  if (error) {
    return <ErrorPage msg={error} />;
  }

  return (
    <>
      <header className="p-3 items-start gap-2 flex flex-col">
        <Input
          className="rounded-full dark:bg-dre-bg-light/5 px-6 py-4 border-none md:w-1/2"
          placeholder="Search the Archive..."
        />

        <div className="flex mx-3">
          <p>Filter</p>
          <p>Sort By</p>
        </div>
      </header>
      <main className="space-y-3 mx-8">
        {data!.map((dispute) => (
          <SearchResult key={dispute.id} {...dispute} />
        ))}
      </main>
      <footer>
        <Pagination>
          <PaginationContent>
            {params.offset > 0 && (
              <>
                <PaginationItem>
                  <PaginationPrevious href={pager(params, params.offset - 1)}>
                    Previous
                  </PaginationPrevious>
                </PaginationItem>
                <p>Page {params.offset + 1}</p>{" "}
              </>
            )}
            {params.offset < total && (
              <PaginationItem>
                <PaginationNext href={pager(params, params.offset + 1)}>Next</PaginationNext>
              </PaginationItem>
            )}
          </PaginationContent>
        </Pagination>
      </footer>
    </>
  );
}
