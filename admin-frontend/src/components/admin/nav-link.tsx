import Link from "next/link";
import { usePathname } from "next/navigation";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";

export default function NavLink({
  label,
  href,
  icon,
  expanded = false,
}: {
  label: string;
  href: string;
  expanded?: boolean;
  icon: ReactNode;
}) {
  const pathname = usePathname();
  const className = cn(
    "p-3 rounded-xl transition-all",
    href == pathname &&
      "hover:shadow-lg hover:shadow-primary-600/35 text-primary-500 dark:text-secondary-500 border border-primary-500/35 bg-primary-500/20",
    expanded && "w-full",
  );
  return (
    <Button asChild className={className} variant="ghost">
      <Link href={href} title={label}>
        {icon}
        {expanded && <span className="ml-3">{label}</span>}
      </Link>
    </Button>
  );
}
