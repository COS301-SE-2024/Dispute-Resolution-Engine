"use client";

import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ChevronRightIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { usePathname } from "next/navigation";

export function DisputeLink({ href, children }: { href: string; children: React.ReactNode }) {
  const pathname = usePathname();
  const c = cn(buttonVariants({ variant: pathname == href ? "default" : "ghost" }), "flex");
  return (
    <Link href={href} className={c}>
      <span className="grow truncate">{children}</span>
      <ChevronRightIcon />
    </Link>
  );
}
