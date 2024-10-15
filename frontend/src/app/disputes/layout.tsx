import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Metadata } from "next";
import ClientSearch from "./clientSearch";
import { MenuIcon } from "lucide-react";

export const metadata: Metadata = {
  title: "DRE - Disputes",
};

export default function DisputeRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-[auto_1fr] h-full">
      <div className="flex-col gap-4 h-full p-4 flex row-start-2 md:row-start-1 border-t md:border-r border-dre-200">
        <ClientSearch />

        <Button asChild className="mt-auto">
          <Link href="/disputes/create" className="w-full">
            + Create
          </Link>
        </Button>
      </div>
      <div className="overflow-y-auto">{children}</div>
      {/* <Separator orientation="vertical" /> */}
    </div>
  );
}
