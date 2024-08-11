"use client";

import { Badge } from "@/components/ui/badge";
import { buttonVariants } from "@/components/ui/button";
import { Role } from "@/lib/interfaces/dispute";
import { cn } from "@/lib/utils";
import { ChevronRight } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

export function DisputeLink({ href, role, title }: { title: string; role: Role; href: string }) {
  const pathname = usePathname();
  const c = cn(buttonVariants({ variant: pathname == href ? "default" : "ghost" }), "flex");
  return (
    <Link href={href} className={c}>
      <span className="grow truncate" title={title}>
        {title}
      </span>
      <div>
        {role == "Complainant" ? (
          <Badge className="ml-2 inline" title="You are a complainant in this case">
            {role.substring(0, 1)}
          </Badge>
        ) : role == "Respondent" ? (
          <Badge
            className="ml-2 inline"
            variant="secondary"
            title="You are a respondant in this case"
          >
            {role.substring(0, 1)}
          </Badge>
        ) : null}
        <ChevronRight size="1rem" className="inline" />
      </div>
    </Link>
  );
}
