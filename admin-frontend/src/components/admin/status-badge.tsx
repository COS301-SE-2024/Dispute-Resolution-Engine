import { cva, VariantProps } from "class-variance-authority";
import { Badge } from "../ui/badge";
import { cn } from "@/lib/utils";
import { FunctionComponent, HTMLAttributes } from "react";
import { ChevronDown } from "lucide-react";
import { TicketStatus } from "@/lib/types/tickets";

const statusVariants = cva("", {
  variants: {
    variant: {
      warning: [
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
      success: [
        "hover:bg-green-500/30 border-green-500/50 bg-green-500/20 text-green-700",
        "dark:hover:bg-green-500/30 dark:border-green-500/70 dark:bg-green-500/20 dark:text-green-400",
      ],
      neutral: [
        "hover:bg-primary-500/30 border-primary-500/50 bg-primary-500/20 text-primary-700",
        "dark:hover:bg-primary-400/30 dark:border-primary-400/70 dark:bg-primary-400/20 dark:text-primary-300",
      ],
    },
  },
  defaultVariants: {
    variant: "neutral",
  },
});

export type StatusVariant = Exclude<
  VariantProps<typeof statusVariants>["variant"],
  null | undefined
>;

export interface StatusBadgeProps
  extends HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof statusVariants> {
  dropdown?: boolean;
}

export function StatusBadge({ children, variant, dropdown = false }: StatusBadgeProps) {
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

export type MappedStatusProps<K extends string> = Omit<StatusBadgeProps, "variant"> & {
  variant: K;
};

/**
 * Utility function that maps a string union to a status color
 * @param mapping  The mapping from the string union to the status color
 * @returns The modified component
 */
function mapStatus<K extends string>(
  mapping: Record<K, StatusVariant>
): FunctionComponent<MappedStatusProps<K>> {
  const component = ({ variant, ...props }: MappedStatusProps<K>) => (
    <StatusBadge variant={mapping[variant]} {...props} />
  );
  component.displayName = "StatusBadge";
  return component;
}

export const TicketStatusBadge = mapStatus<TicketStatus>({
  Open: "neutral",
  Closed: "inactive",
  Solved: "success",
  "On Hold": "warning",
});
