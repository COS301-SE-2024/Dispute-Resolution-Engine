import * as React from "react";

import { cn } from "@/lib/utils";

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex w-full rounded-md border border-primary-500/30 bg-surface-light-100 px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 disabled:text-black/50 disabled:dark:text-white/50 disabled:cursor-not-allowed dark:border-primary-500/30 dark:bg-surface-dark-900 dark:placeholder:text-slate-400 focus-visible:ring-primary-500 focus-visible:bg-white dark:focus-visible:bg-surface-dark-800 transition-colors",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Input.displayName = "Input";

export { Input };
