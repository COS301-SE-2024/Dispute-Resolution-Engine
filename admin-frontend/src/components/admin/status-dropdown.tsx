"use client";
import { cva, VariantProps } from "class-variance-authority";
import { Badge } from "../ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { cn } from "@/lib/utils";
import { ChevronDown } from "lucide-react";
import { ReactNode } from "react";
import { DisputeStatus } from "@/lib/types/dispute";

const statusVariants = cva("", {
  variants: {
    variant: {
      waiting: [
        "hover:bg-yellow-500/30 border-yellow-500/50 bg-yellow-500/20 text-yellow-500",
        "dark:hover:bg-yellow-500/30 dark:border-yellow-500/70 dark:bg-yellow-500/20 dark:text-yellow-400",
      ],
      error: [
        "hover:bg-red-500/30 border-red-500/40 bg-red-500/20 text-red-700",
        "dark:hover:bg-red-500/30 dark:border-red-500/70 dark:bg-red-500/20 dark:text-red-400",
      ],
      inactive: [
        "hover:bg-slate-500/30 border-slate-500/40 bg-slate-500/20 text-slate-700",
        "dark:hover:bg-slate-500/30 dark:border-slate-500/70 dark:bg-slate-500/20 dark:text-slate-300",
      ],
      active: [
        "hover:bg-green-500/30 border-green-500/50 bg-green-500/20 text-green-700",
        "dark:hover:bg-green-500/30 dark:border-green-500/70 dark:bg-green-500/20 dark:text-green-400",
      ],
    },
  },
});

export interface StatusProps extends VariantProps<typeof statusVariants> {
  dropdown?: boolean;
  children?: ReactNode;
}

export function StatusBadge({ children, variant, dropdown = false }: StatusProps) {
  return (
    <Badge
      className={cn(
        dropdown && "pl-1",
        statusVariants({
          variant,
        })
      )}
    >
      {dropdown && <ChevronDown size="1rem" />}
      {children}
    </Badge>
  );
}

export function StatusDropdown({
  onSelect = () => {},
  children,
  disabled,
}: {
  onSelect?: (status: DisputeStatus) => void;
  children: ReactNode;
  disabled?: boolean;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger disabled={disabled} className="disabled:opacity-50">
        {children}
      </DropdownMenuTrigger>
      <DropdownMenuContent className="rounded-md">
        <DropdownMenuItem onSelect={() => onSelect("Awaiting respondent")}>
          <StatusBadge variant="waiting">Awaiting respondent</StatusBadge>
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => onSelect("Active")}>
          <StatusBadge variant="active">Active</StatusBadge>
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => onSelect("Review")}>
          <StatusBadge variant="error">Review</StatusBadge>
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => onSelect("Settled")}>
          <StatusBadge variant="inactive">Settled</StatusBadge>
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => onSelect("Refused")}>
          <StatusBadge variant="inactive">Refused</StatusBadge>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
