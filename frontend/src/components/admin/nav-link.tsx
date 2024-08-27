import Link from "next/link";

import { Button } from "@/components/ui/button";

import { cn } from "@/lib/utils";

export default function NavLink({
  label,
  href,
  icon,
  expanded = false,
  active = false,
}: {
  label: string;
  href: string;
  expanded?: boolean;
  active?: boolean;
  icon: ReactNode;
}) {
  const className = cn(
    "w-auto p-3 h-auto rounded-xl",
    expanded && "w-full",
    active && "text-secondary-500 border border-primary-500/35 bg-primary-500/20",
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
