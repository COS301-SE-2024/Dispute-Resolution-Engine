import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ExternalLink } from "lucide-react";
import Link from "next/link";

function SearchResult() {
  return (
    <Card className="border-none">
      <CardHeader>
        <Link href="/archive/id">
          <CardTitle>Dispute Title</CardTitle>
        </Link>
      </CardHeader>
      <CardContent asChild className="dark:text-white/50">
        <p>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Molestias quo eveniet, veniam
          facere voluptatem est...
        </p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <p>Status: Resolved</p>
        <Button>
          <ExternalLink size="1rem" className="mr-2" />
          Read More
        </Button>
      </CardFooter>
    </Card>
  );
}

export default function ArchiveSearch() {
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
        {[...Array(10).keys()].map((i) => (
          <SearchResult key={i} />
        ))}
      </main>
    </>
  );
}
