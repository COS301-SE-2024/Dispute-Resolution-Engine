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

import { TICKET_STATUS, TicketStatus } from "@/lib/types/tickets";
import { TicketStatusBadge } from "./status-badge";
import { DISPUTE_STATUS, DisputeStatus } from "@/lib/types/dispute";

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

type Variant = VariantProps<typeof statusVariants>["variant"];

const STATUS_VARIANTS: {
  status: DisputeStatus;
  variant: Variant;
}[] = [
  { status: "Awaiting Respondant", variant: "waiting" },
  { status: "Active", variant: "active" },
  { status: "Review", variant: "error" },
  { status: "Refused", variant: "error" },
  { status: "Appeal", variant: "inactive" },
  { status: "Settled", variant: "inactive" },
  { status: "Withdrawn", variant: "inactive" },
  { status: "Transfer", variant: "inactive" },
  { status: "Other", variant: "inactive" },
] as const;

export interface StatusProps {
  value: DisputeStatus;
  dropdown?: boolean;
}

export function StatusBadge({ value, dropdown = false }: StatusProps) {
  const variant = STATUS_VARIANTS.find(({ status }) => status == value)?.variant;
  if (!variant) {
    throw new Error("Unhandled dispute status: " + value);
  }

  return dropdown ? (
    <Badge
      className={cn(
        "pl-1",
        statusVariants({
          variant,
        })
      )}
      asChild
    >
      <button>
        {dropdown && <ChevronDown size="1rem" />}
        {value}
      </button>
    </Badge>
  ) : (
    <Badge className={statusVariants({ variant })}>
      {dropdown && <ChevronDown size="1rem" />}
      {value}
    </Badge>
  );
}

export function StatusDropdown({
  onSelect = () => {},
  initialValue,
  children,
  disabled,
}: {
  onSelect?: (status: DisputeStatus) => void;
  initialValue?: DisputeStatus;
  children: ReactNode;
  disabled?: boolean;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger disabled={disabled} className="disabled:opacity-50">
        {children}
      </DropdownMenuTrigger>
      <DropdownMenuContent className="rounded-md">
        {DISPUTE_STATUS.filter((status) => status !== initialValue).map((value) => (
          <DropdownMenuItem key={value} onSelect={() => onSelect(value)}>
            <StatusBadge value={value} />
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export function TicketStatusDropdown({
  onSelect = () => {},
  children,
  disabled,
}: {
  onSelect?: (status: TicketStatus) => void;
  children: ReactNode;
  disabled?: boolean;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger disabled={disabled} className="disabled:opacity-50">
        {children}
      </DropdownMenuTrigger>
      <DropdownMenuContent className="rounded-md">
        {TICKET_STATUS.map((status) => (
          <DropdownMenuItem key={status} onSelect={() => onSelect(status)}>
            <TicketStatusBadge variant={status}>{status}</TicketStatusBadge>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
