import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { searchArchive } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import { ExternalLink } from "lucide-react";
import Link from "next/link";
import { redirect } from "next/navigation";

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

interface SearchParams {
  q?: string;
}

export default async function ArchiveSearch({ searchParams }: { searchParams: SearchParams }) {
  if (!searchParams.q) {
    redirect("/archive");
  }

  const { data, error } = await searchArchive({
    search: searchParams.q,
  });

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-5 w-2/3 mx-auto">
        <h1 className="text-4xl font-bold tracking-wide">Oops, something went wrong :(</h1>
        <p className="dark:text-white/50">{error}</p>
      </div>
    );
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
    </>
  );
}
