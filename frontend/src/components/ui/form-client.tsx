"use client";

import { cn } from "@/lib/utils";
import { HTMLAttributes, ReactNode } from "react";
import { Label } from "./label";
import { FieldValues, useFormContext } from "react-hook-form";

export type FormMessageProps<T> = Omit<HTMLAttributes<HTMLParagraphElement>, "children"> & {
  name?: keyof T;
};

export function FormMessage<T extends FieldValues>({
  name,
  className,
  ...props
}: FormMessageProps<T>) {
  const {
    formState: { errors },
  } = useFormContext<T>();
  const error = errors[name ?? "root"]?.message as string | undefined;
  if (!error) {
    return null;
  }

  return (
    <p
      role="alert"
      className={cn("text-xs font-medium text-red-500 dark:text-red-500", className)}
      {...props}
    >
      {error}
    </p>
  );
}

export function FormField<T extends FieldValues>({
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
  const {
    formState: { errors },
  } = useFormContext<T>();

  return (
    <div className={className}>
      <Label htmlFor={id} className={cn(errors[name] ? "text-red-500" : undefined)}>
        {label}
      </Label>
      {children}
      <Message name={name} />
    </div>
  );
}
