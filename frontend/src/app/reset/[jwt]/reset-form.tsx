"use client";

import { useFormStatus } from "react-dom";
import { resetPassword } from "@/lib/actions/auth";
import { Label } from "@/components/ui/label";
import { ResetPassError, ResetPassData } from "@/lib/schema/auth";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ReactNode, forwardRef, useContext } from "react";
import { Result } from "@/lib/types";
import { createFormContext } from "@/components/ui/form-server";

const [ResetContext, ResetForm] = createFormContext<Result<string, ResetPassError>>(
  "ResetForm",
  resetPassword
);
export { ResetForm };

const FormMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <p
        ref={ref}
        className={cn("text-xs font-medium text-red-500 dark:text-red-500", className)}
        {...props}
      >
        {children}
      </p>
    );
  }
);
FormMessage.displayName = "FormMessage";

const ResetMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  (props, ref) => {
    const state = useContext(ResetContext);
    const error = state?.error && state.error._errors?.at(0);
    return (
      <FormMessage {...props} ref={ref}>
        {error}
      </FormMessage>
    );
  }
);
ResetMessage.displayName = "ResetMessage";
export { ResetMessage };

export function ResetButton() {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      Reset Password
    </Button>
  );
}

export function ResetField({
  name,
  label,
  children,
  className = "",
}: {
  name: keyof ResetPassData;
  label: string;
  children: ReactNode;
  className?: string;
}) {
  const state = useContext(ResetContext);
  const error = state?.error && state.error[name]?._errors?.at(0);
  return (
    <div className={className}>
      <Label htmlFor={name}>{label}</Label>
      {children}
      {error && <FormMessage>{error}</FormMessage>}
    </div>
  );
}
