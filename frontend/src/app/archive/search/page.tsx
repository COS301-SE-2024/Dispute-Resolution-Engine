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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useRef } from "react";
import SearchFilters from "./filters";
import { Badge } from "@/components/ui/badge";

const searchSchema = z.object({
  q: z.string().optional(),
  offset: z.coerce.number().default(0),
  order: z
    .enum(["asc", "desc", ""])
    .transform((e) => (e.length == 0 ? undefined : e))
    .optional(),
  sort: z
    .enum(["title", "date_filed", "date_resolved", "date_filed", "time_taken", ""])
    .transform((e) => (e.length == 0 ? undefined : e))
    .optional(),
});

export type SearchParams = z.infer<typeof searchSchema>;

function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <li>
      <div className="flex items-center gap-5 mb-3">
        <Link href={`/archive/${props.id}`}>
          <h3 className="hover:underline font-semibold text-lg">{props.title}</h3>
        </Link>
        <div className="space-x-1">
          {props.category.map((cat) => (
            <Badge key={cat}>{cat}</Badge>
          ))}
        </div>
      </div>
      <p className="dark:text-white/50">{props.summary}</p>
    </li>
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
    offset: params.offset * 10,
    limit: 10,
  });

  // TODO: replace this with response information
  const total = 10;

  if (error) {
    return <ErrorPage msg={error} />;
  }

  return (
    <>
      <form className="p-3 items-start gap-2 flex flex-col">
        <Input
          defaultValue={params.q}
          name="q"
          className="rounded-full dark:bg-dre-bg-light/5 px-6 py-4 border-none md:w-1/2"
          placeholder="Search the Archive..."
        />
        <SearchFilters params={params} />
      </form>
      <main className="mx-20 grid grid-cols-2">
        <ol className="space-y-5">
          {data!.map((dispute) => (
            <SearchResult key={dispute.id} {...dispute} />
          ))}
        </ol>
      </main>
      <footer className="my-10">
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
