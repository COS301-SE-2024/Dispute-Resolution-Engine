import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { ExternalLink } from "lucide-react";
import Link from "next/link";

function ArchivedDispute() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Dispute Title</CardTitle>
      </CardHeader>
      <CardContent asChild className="dark:text-white/50">
        <p>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Molestias quo eveniet, veniam
          facere voluptatem est...
        </p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <p>Status: Resolved</p>
        <Button asChild>
          <Link href="/id">
            <ExternalLink size="1rem" className="mr-2" />
            Read More
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
export default function Archive() {
  return (
    <div className="flex flex-col items-center justify-center h-full gap-5 w-2/3 mx-auto">
      <header className="mx-auto w-fit text-center">
        <h1 className="text-4xl font-bold tracking-wide">Archive</h1>
        <p className="dark:text-white/50">Explore our previously handled cases</p>
      </header>
      <main className="w-2/3">
        <form action="/archive/search">
          <Input
            name="q"
            className="rounded-full dark:bg-dre-bg-light/5 px-8 py-6 border-none"
            placeholder="Search the Archive..."
          />
        </form>
      </main>
      <footer>
        <h2 className="text-2xl font-semibold mb-4">Resolved Disputes</h2>
        <div className="flex flex-col md:grid md:grid-cols-3 gap-4">
          {[...Array(3).keys()].map((i) => (
            <ArchivedDispute key={i} />
          ))}
        </div>
      </footer>
    </div>
  );
}
