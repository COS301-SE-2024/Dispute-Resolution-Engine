"use client";

import { useFormStatus } from "react-dom";
import { signup } from "@/lib/actions/auth";
import { Label } from "@/components/ui/label";
import { SignupData, SignupError } from "@/lib/schema/auth";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ReactNode, forwardRef, useContext } from "react";
import { Result } from "@/lib/types";
import { createFormContext } from "@/components/ui/form-server";

const [SignupContext, SignupForm] = createFormContext<Result<string, SignupError>>(
  "SignupForm",
  signup
);
export { SignupForm };

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

const SignupMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  (props, ref) => {
    const state = useContext(SignupContext);
    const error = state?.error && state.error._errors?.at(0);
    return (
      <FormMessage {...props} ref={ref}>
        {error}
      </FormMessage>
    );
  }
);
SignupMessage.displayName = "SignupMessage";
export { SignupMessage };

export function SignupButton() {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      Signup
    </Button>
  );
}

export function SignupField({
  name,
  label,
  children,
  className = "",
}: {
  name: keyof SignupData;
  label: string;
  children: ReactNode;
  className?: string;
}) {
  const state = useContext(SignupContext);
  const error = state?.error && state.error[name]?._errors?.at(0);
  return (
    <div className={className}>
      <Label htmlFor={name}>{label}</Label>
      {children}
      {error && <FormMessage>{error}</FormMessage>}
    </div>
  );
}
