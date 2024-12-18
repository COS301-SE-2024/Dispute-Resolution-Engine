import Loader from "@/components/Loader";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { fetchArchiveHighlights } from "@/lib/api/archive";
import { ArchivedDisputeSummary, ArchiveSearchResponse } from "@/lib/interfaces/archive";
import { ExternalLink } from "lucide-react";
import { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "DRE - Archive",
};

function ArchivedDispute(props: ArchivedDisputeSummary) {
  return (
    <Card className="flex flex-col max-w-sm">
      <CardHeader className="flex flex-row items-center justify-between flex-wrap">
        <div>
          <CardTitle>{props.title}</CardTitle>
          <p>{props.date_resolved}</p>
        </div>
        <ul className="flex gap-1 flex-wrap">
          {props.category.map((cat, i) => (
            <Badge key={i}>{cat}</Badge>
          ))}
        </ul>
      </CardHeader>
      <CardContent asChild className="dark:text-white/50 grow truncate">
        <p>{props.summary}</p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button asChild>
          <Link href={`/archive/${props.id}`}>
            <ExternalLink size="1rem" className="mr-2" />
            Read More
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

export default async function Archive() {
  const { data, error } = await fetchArchiveHighlights(3);
  if (error) {
    return <h1>{error}</h1>;
  }
  return (
    <div className="flex flex-col items-center justify-center h-full gap-5 w-2/3 mx-auto">
      <header className="mx-auto w-fit text-center">
        <h1 className="text-6xl font-bold tracking-wide">Archive</h1>
        <p className="dark:text-white/50">Explore our previously handled cases</p>
      </header>
      <main className="w-2/3">
        <form action="/archive/search" className="flex flex-col items-center gap-2">
          <Input
            name="q"
            className="rounded-full dark:bg-dre-bg-light/5 px-6 py-4 border-none"
            placeholder="Search the Archive..."
          />
          <div className="flex gap-2">
            <Button type="submit">Search</Button>
            <Button asChild>
              <Link href="/archive/browse">
                <ExternalLink size="1rem" className="mr-2" />
                Browse
              </Link>
            </Button>
          </div>
        </form>
      </main>
      <footer>
        <h2 className="text-2xl font-semibold mb-4">Resolved Disputes</h2>
        <div className="flex flex-col md:grid md:grid-cols-3 gap-4">
          {data!.archives.map((props, i) => (
            <ArchivedDispute key={i} {...props} />
          ))}
        </div>
      </footer>
    </div>
  );
}
