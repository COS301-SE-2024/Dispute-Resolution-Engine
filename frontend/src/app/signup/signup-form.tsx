"use client";

import { useFormState, useFormStatus } from "react-dom";
import { signup } from "../lib/auth/actions";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { SignupData, SignupError } from "../lib/auth/types";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { forwardRef, useId } from "react";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

function TextField({
  name,
  label,
  type,
  state,
}: {
  name: keyof SignupData;
  label: string;
  type?: "text" | "password";
  state?: SignupError;
}) {
  const error = state && state[name]?._errors?.at(0);
  return (
    <>
      <Label htmlFor={name}>{label}</Label>
      <Input type={type} id={name} name={name} placeholder={label} />
      {error && <FormMessage>{error}</FormMessage>}
    </>
  );
}

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

function SignupButton() {
  const { pending } = useFormStatus();
  return (
    <Button disabled={pending} type="submit">
      Create
    </Button>
  );
}

export default function SignupForm() {
  const [state, formAction] = useFormState(signup, undefined);

  return (
    <Card asChild className="mx-auto md:my-3 lg:w-1/2 md:w-3/4">
      <form action={formAction}>
        <CardHeader>
          <CardTitle>Create an Account</CardTitle>
        </CardHeader>
        <CardContent asChild>
          <Tabs defaultValue="profile">
            <TabsList>
              <TabsTrigger value="profile">Profile</TabsTrigger>
              <TabsTrigger value="address">Address</TabsTrigger>
            </TabsList>
            <TabsContent value="profile" forceMount className="data-[state=inactive]:hidden">
              <TextField state={state?.error} name="firstName" label="First Name" type="text" />
              <TextField state={state?.error} name="lastName" label="Last Name" type="text" />
              <TextField state={state?.error} name="email" label="Email" type="text" />
              <TextField state={state?.error} name="password" label="Password" type="password" />
              <TextField
                state={state?.error}
                name="passwordConfirm"
                label="Confirm Password"
                type="password"
              />
            </TabsContent>
            <TabsContent value="address" forceMount className="data-[state=inactive]:hidden">
              <TextField state={state?.error} name="addrCountry" label="Country" type="text" />
              <TextField state={state?.error} name="addrProvince" label="Province" type="text" />
              <TextField state={state?.error} name="addrCity" label="City" type="text" />
              <TextField state={state?.error} name="addrStreet" label="Street 1" type="text" />
              <TextField state={state?.error} name="addrStreet2" label="Street 2" type="text" />
              <TextField state={state?.error} name="addrStreet3" label="Street 3" type="text" />
            </TabsContent>
          </Tabs>
        </CardContent>
        <CardFooter>
          <SignupButton />
          <p role="alert">{state?.data}</p>
        </CardFooter>
      </form>
    </Card>
  );
}
