import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import { ChevronRightIcon } from "@radix-ui/react-icons";
import { fetchDisputes } from "../lib/dispute";

function DisputeLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <Button variant="ghost" className="w-full grow text-left">
      <Link href={href} className="inline-flex w-full">
        <span className="grow">{children}</span>
        <ChevronRightIcon />
      </Link>
    </Button>
  );
}

export default async function DisputeRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const data = await fetchDisputes("cook");

  return (
    <div className="flex items-stretch h-full lg:w-3/4 mx-auto">
      <div className="w-56 flex flex-col p-2 gap-4">
        <Input placeholder="Search" />
        <nav>
          <ul>
            {typeof data != "string" ? data.map(d => (
            <li key={d.id}>
              <DisputeLink href={`/disputes/${d.id}`}>{d.title}</DisputeLink>
            </li>
            )) : <li>{data}</li>}
          </ul>
        </nav>

        <Button className="mt-auto" asChild variant="outline">
          <Link href="/disputes/create" className="w-full">
            + Create
          </Link>
        </Button>
      </div>
      <Separator orientation="vertical" />
      {children}
    </div>
  );
}
