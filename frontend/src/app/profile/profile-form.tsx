"use client";

import { ProfileData, ProfileError } from "@/lib/schema/profile";
import { cn } from "@/lib/utils";
import { ReactNode, forwardRef, useContext } from "react";
import { Result } from "@/lib/types";
import { createFormContext } from "@/components/ui/form-server";
import { updateProfile } from "@/lib/actions/profile";
import { useFormStatus } from "react-dom";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";

const [ProfileUpdateContext, ProfileUpdateForm] = createFormContext<Result<string, ProfileError>>(
  "ProfileUpdateForm",
  updateProfile
);
export { ProfileUpdateForm };

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

const ProfileUpdateMessage = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  (props, ref) => {
    const state = useContext(ProfileUpdateContext);
    const error = state?.error && state.error._errors?.at(0);
    return (
      <FormMessage {...props} ref={ref}>
        {error}
      </FormMessage>
    );
  }
);
ProfileUpdateMessage.displayName = "ProfileUpdateMessage";
export { ProfileUpdateMessage };

export function ProfileUpdateField({
  name,
  label,
  children,
  className = "",
  id = name
}: {
  name: keyof ProfileData;
  label: string;
  children: ReactNode;
  className?: string;
  id?: string;
}) {
  const state = useContext(ProfileUpdateContext);
  const error = state?.error && state.error[name]?._errors?.at(0);
  return (
    <div className={className}>
      <Label htmlFor={id}>{label}</Label>
      {children}
      {error && <FormMessage>{error}</FormMessage>}
    </div>
  );
}

export function ProfileUpdateButton() {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      Save
    </Button>
  );
}
