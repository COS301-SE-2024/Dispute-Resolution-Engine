import { cn } from "@/lib/utils";
import { forwardRef, HTMLAttributes, ReactNode } from "react";

export const Root = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      {...props}
      className={cn("grid grid-rows-[auto_1fr] overflow-y-hidden", className)}
    />
  )
);
Root.displayName = "Root";

export const Header = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <header ref={ref} {...props} className={cn("p-3 border-b border-dre-200/30", className)} />
  )
);
Header.displayName = "Header";

export const Content = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <header ref={ref} {...props} className={cn("p-4 py-6 overflow-y-auto", className)} />
  )
);
Content.displayName = "Content";
