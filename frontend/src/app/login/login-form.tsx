"use client";

import { useFormStatus } from "react-dom";
import { login } from "@/lib/actions/auth";
import { Label } from "@/components/ui/label";
import { LoginData, LoginError, loginSchema } from "@/lib/schema/auth";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { ReactNode, forwardRef, useContext } from "react";
import { Result } from "@/lib/types";
import { createFormContext } from "@/components/ui/form-server";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const [LoginContext, LoginForm] = createFormContext<Result<string, LoginError>>("LoginForm", login);
export { LoginForm };

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

const LoginMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  (props, ref) => {
    const state = useContext(LoginContext);
    const error = state?.error && state.error._errors?.at(0);
    return (
      <FormMessage {...props} ref={ref}>
        {error}
      </FormMessage>
    );
  }
);
LoginMessage.displayName = "LoginMessage";
export { LoginMessage };

export function LoginButton() {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      Login
    </Button>
  );
}

export function LoginField({
  name,
  label,
  children,
}: {
  name: keyof LoginData;
  label: string;
  children: ReactNode;
}) {
  const state = useContext(LoginContext);
  const error = state?.error && state.error[name]?._errors?.at(0);
  return (
    <div className="space-y-1">
      <Label htmlFor={name}>{label}</Label>
      {children}
      {error && <FormMessage>{error}</FormMessage>}
    </div>
  );
}
