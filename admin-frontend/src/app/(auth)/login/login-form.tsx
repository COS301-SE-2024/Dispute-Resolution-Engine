"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { login } from "@/lib/api/auth";
import { useId } from "react";
import { useFormState, useFormStatus } from "react-dom";

export default function LoginForm() {
  const emailId = useId();
  const passId = useId();

  const formId = useId();

  const [state, formAction] = useFormState(login, undefined);
  const { pending } = useFormStatus();

  return (
    <Card className="md:mx-auto md:max-w-xl mt-5 mx-2 sm:mx-5">
      <CardHeader>
        <CardTitle>Admin Login</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <form id={formId} action={formAction} className="space-y-4">
          <div>
            <Label htmlFor={emailId}>Email</Label>
            <Input id={emailId} name="email" autoComplete="email" placeholder="Email" />
          </div>

          <div>
            <Label htmlFor={passId}>Password</Label>
            <Input
              id={passId}
              name="password"
              autoComplete="current-password"
              placeholder="Password"
              type="password"
            />
          </div>
        </form>
      </CardContent>
      <CardFooter className="mt-auto flex justify-between">
        {state?.error ? (
          <p className="text-red-500 text-sm" role="alert">
            {state.error}
          </p>
        ) : (
          <div></div>
        )}
        <Button disabled={pending} form={formId} type="submit" className="ml-auto">
          Login
        </Button>
      </CardFooter>
    </Card>
  );
}
