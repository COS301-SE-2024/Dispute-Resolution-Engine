import * as React from "react";

import { cn } from "@/lib/utils";
import { forwardSlot } from "../util";

const Card = forwardSlot<HTMLDivElement>(
  "Card",
  "div",
  "bg-surface-light-50 dark:bg-surface-dark-800 rounded-lg dark:border dark:border-primary-500/10 max-w-5xl mx-auto shadow-md"
);

const CardHeader = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn("flex flex-col space-y-1.5 px-6 pt-6", className)} {...props} />
  )
);
CardHeader.displayName = "CardHeader";

const CardTitle = React.forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => (
    <h3
      ref={ref}
      className={cn("text-2xl font-semibold leading-none tracking-tight", className)}
      {...props}
    />
  )
);
CardTitle.displayName = "CardTitle";

const CardDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <p ref={ref} className={cn("text-sm text-slate-500 dark:text-slate-400", className)} {...props} />
));
CardDescription.displayName = "CardDescription";

const CardContent = forwardSlot<HTMLDivElement>("CardContent", "div", "p-6");
const CardFooter = forwardSlot<HTMLDivElement>("CardFooter", "div", "flex items-center p-6 pt-0");

export { Card, CardHeader, CardFooter, CardTitle, CardDescription, CardContent };
