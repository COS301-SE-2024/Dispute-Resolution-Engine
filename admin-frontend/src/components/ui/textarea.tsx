import * as React from "react";

import { cn } from "@/lib/utils";

export interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {}

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, ...props }, ref) => {
    return (
      <textarea
        className={cn(
          "flex min-h-[80px] w-full rounded-md border border-primary-500/30 bg-surface-light-100 px-3 py-2 text-sm placeholder:text-slate-500 focus-visible:bg-white focus-visible:outline-none focus-visible:ring-2 ring-primary-500 disabled:cursor-not-allowed disabled:text-black/50 dark:disabled:text-white/50 dark:border-primary-500/30 dark:bg-surface-dark-900 dark:placeholder:text-slate-400 dark:focus-visible:bg-surface-dark-800 transition-colors",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Textarea.displayName = "Textarea";

export { Textarea };
