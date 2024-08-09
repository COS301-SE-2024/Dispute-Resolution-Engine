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
import ClientSearch from "./clientSearch";

export const metadata: Metadata = {
  title: "DRE - Disputes",
};



export default function DisputeRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex items-stretch h-full lg:w-3/4 mx-auto">
      <div className="flex shrink-0 flex-col -2 gap-4">
        <ClientSearch></ClientSearch>

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
