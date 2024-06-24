"use client";

import { useFormStatus } from "react-dom";
import { verify } from "@/lib/actions/auth";
import { Label } from "@/components/ui/label";
import { VerifyData, VerifyError } from "@/lib/schema/auth";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ReactNode, forwardRef, useContext } from "react";
import { Result } from "@/lib/types";
import { createFormContext } from "@/components/ui/form-server";

const [VerifyContext, VerifyForm] = createFormContext<Result<string, VerifyError>>(
  "VerifyForm",
  verify
);
export { VerifyForm };

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

const VerifyMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  (props, ref) => {
    const state = useContext(VerifyContext);
    const error = state?.error && state.error._errors?.at(0);
    return (
      <FormMessage {...props} ref={ref}>
        {error}
      </FormMessage>
    );
  }
);
VerifyMessage.displayName = "VerifyMessage";
export { VerifyMessage };

export function VerifyButton() {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      Verify
    </Button>
  );
}

export function VerifyField({
  name,
  label,
  children,
  className = "",
}: {
  name: keyof VerifyData;
  label: string;
  children: ReactNode;
  className?: string;
}) {
  const state = useContext(VerifyContext);
  const error = state?.error && state.error[name]?._errors?.at(0);
  console.log(name, error);
  return (
    <div className={className}>
      <Label htmlFor={name}>{label}</Label>
      {children}
      {error && <FormMessage>{error}</FormMessage>}
      <FormMessage>{error}</FormMessage>
    </div>
  );
}
