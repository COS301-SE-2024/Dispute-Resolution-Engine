import Link from "next/link";
import { Button, buttonVariants } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import { ChevronRightIcon } from "@radix-ui/react-icons";
import { getDisputeList } from "../../lib/api/dispute";
import { Suspense } from "react";
import Loader from "@/components/Loader";
import { Metadata } from "next";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

export const metadata: Metadata = {
  title: "DRE - Disputes",
};

function DisputeLink({ href, children }: { href: string; children: React.ReactNode }) {
  const c = cn(buttonVariants({ variant: "ghost" }), "flex");
  return (
    <Link href={href} className={c}>
      <span className="grow truncate">{children}</span>
      <ChevronRightIcon />
    </Link>
  );
}

async function DisputeList() {
  const data = await getDisputeList();

  return (
    <ul>
      {data.data ? (
        data.data.map((d) => (
          <li key={d.id}>
            <DisputeLink href={`/disputes/${d.id}`}>
              {d.title}
              {d.role == "Complainant" ? (
                <Badge className="ml-2">{d.role.substring(0, 1)}</Badge>
              ) : d.role == "Respondant" ? (
                <Badge className="ml-2" variant="secondary">
                  {d.role.substring(0, 1)}
                </Badge>
              ) : null}
            </DisputeLink>
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
      <div className="flex shrink-0 flex-col p-2 gap-4">
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
