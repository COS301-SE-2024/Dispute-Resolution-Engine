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
    <div className="grid grid-cols-1 md:grid-cols-[auto_1fr] overflow-y-hidden">
      <ClientSearch />
      {children}
    </div>
  );
}
