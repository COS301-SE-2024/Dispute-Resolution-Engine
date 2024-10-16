import { cn } from "@/lib/utils";
import { Slot } from "@radix-ui/react-slot";
import React, { HTMLAttributes } from "react";

/**
 * Utiltity function that creates a slottable forwardRef component
 * @param displayName  The display name of the component
 * @param elem The actual tag name to forward ref
 * @param defaultStyles The default styles of the element
 * @returns The created forwardRef
 */
export function forwardSlot<T extends HTMLDivElement>(
  displayName: string,
  elem: keyof HTMLElementTagNameMap,
  defaultStyles: string
) {
  const Comp = React.forwardRef<
    T,
    HTMLAttributes<T> & {
      asChild?: boolean;
    }
  >(({ asChild, className, ...props }, ref) => {
    const Comp = asChild ? Slot : (elem as string);
    return <Comp ref={ref} className={cn(defaultStyles, className)} {...props} />;
  });
  Comp.displayName = displayName;
  return Comp;
}
