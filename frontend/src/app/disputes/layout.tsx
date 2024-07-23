import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import { ChevronRightIcon } from "@radix-ui/react-icons";
import { getDisputeList } from "../../lib/api/dispute";
import { Suspense } from "react";
import Loader from "@/components/Loader";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "DRE - Disputes",
};

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

async function DisputeList() {
  const data = await getDisputeList();

  return (
    <ul>
      {data.data ? (
        data.data.map((d) => (
          <li className="w-auto overflow-x-auto" key={d.id}>
            <DisputeLink href={`/disputes/${d.id}`}>{d.title}</DisputeLink>
          </li>
        ))
      ) : (
        <li>{data.error}</li>
      )}
    </ul>
  );
}

export default function DisputeRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex items-stretch h-full lg:w-3/4 mx-auto">
      <div className="w-56 flex shrink-0 flex-col p-2 gap-4">
        <Input placeholder="Search" />
        <nav>
          <Suspense fallback={<Loader />}>
            <DisputeList />
          </Suspense>
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
