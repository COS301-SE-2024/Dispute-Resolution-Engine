"use client";

import { Result } from "@/lib/types";
import { cn } from "@/lib/utils";
import {
  createContext,
  FormHTMLAttributes,
  forwardRef,
  HTMLAttributes,
  ReactNode,
  useContext,
} from "react";
import { useFormState, useFormStatus } from "react-dom";
import { z } from "zod";
import { Label } from "./label";
import { Button } from "./button";

const FormContext = createContext<any | undefined>(undefined);

type FormState<TIn, TOut = string> = Result<TOut, z.ZodFormattedError<TIn>>;
type FormAction<TIn, TOut = string> = (
  state: Awaited<FormState<TIn> | undefined>,
  data: FormData
) => FormState<TIn> | Promise<FormState<TIn>>;

export type FormProps<TIn> = Omit<FormHTMLAttributes<HTMLFormElement>, "action"> & {
  action: FormAction<TIn>;
};

export function Form<T>({ action, children, ...props }: FormProps<T>) {
  const [state, formAction] = useFormState(action, undefined);

  return (
    <FormContext.Provider value={state}>
      <form action={formAction} {...props}>
        {children}
      </form>
    </FormContext.Provider>
  );
}

const _FormMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
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
_FormMessage.displayName = "FormMessage";

export type FormMessageProps<TIn> = Omit<HTMLAttributes<HTMLParagraphElement>, "children"> & {
  formName?: keyof TIn;
};

export function FormMessage<T>({ formName: key, ...props }: FormMessageProps<T>) {
  const state = useContext(FormContext) as FormState<T> | undefined;

  let error: string | undefined;
  if (key) {
    // z.ZodFormattedError doesn't like accesses using keyof, hence the "as any"
    error = state?.error && (state.error as any)[key]?._errors?.at(0);
  } else {
    error = state?.error && state.error._errors?.at(0);
  }

  if (!error) {
    return undefined;
  }

  return <_FormMessage {...props}>{error}</_FormMessage>;
}

export function FormField<T>({
  id,
  name,
  label,
  children,
  className = "",
}: {
  id: string;
  name: keyof T;
  label: string;
  children: ReactNode;
  className?: string;
}) {
  const Message = FormMessage<T>;

  return (
    <div className={className}>
      <Label htmlFor={id}>{label}</Label>
      {children}
      <Message formName={name} />
    </div>
  );
}

export function FormSubmit({ children }: { children: ReactNode }) {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      {children}
    </Button>
  );
}
