import Link from "next/link";
import { Button, buttonVariants } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import { getDisputeList } from "../../lib/api/dispute";
import { Suspense } from "react";
import Loader from "@/components/Loader";
import { Metadata } from "next";
import { Badge } from "@/components/ui/badge";
import { DisputeLink } from "./link";

export const metadata: Metadata = {
  title: "DRE - Disputes",
};

async function DisputeList() {
  const data = await getDisputeList();

  const err =
    data.error ??
    (data.data.length == 0 ? "You aren't  involved in any disputes. Yay :)" : undefined);

  return (
    <ul>
      {err ? (
        <li className="text-dre-bg-light/50 w-full">{err}</li>
      ) : (
        data.data!.map((d) => (
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

        <Button asChild className="mt-auto">
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
