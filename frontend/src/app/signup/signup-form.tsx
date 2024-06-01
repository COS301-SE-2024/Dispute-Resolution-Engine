"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Form } from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { HTMLAttributes, useId, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { signup } from "@/app/lib/auth";
import TextField from "@/components/form/text-field";
import { sign } from "crypto";

const signupSchema = z
  .object({
    firstName: z.string().min(1, "Required"),
    lastName: z.string().min(1, "Required"),
    email: z.string().min(1, "Required").email("Please enter a valid email"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters long")
      .regex(/\d/gm, "Password must contain at least one digit")
      .regex(/[A-Za-z]/gm, "Password must contain at least one letter")
      .regex(/[^\w\d\s:]/gm, "Password must contain a special character"),
    passwordConfirm: z.string(),
  })
  .superRefine((arg, ctx) => {
    if (arg.password !== arg.passwordConfirm) {
      ctx.addIssue({
        code: "custom",
        message: "The passwords did not match",
        path: ["passwordConfirm"],
      });
    }
  });

export type SignupData = z.infer<typeof signupSchema>;

const SignupField = TextField<SignupData>;

export default function SignupForm(props: HTMLAttributes<HTMLFormElement>) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const form = useForm<SignupData>({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      password: "",
      passwordConfirm: "",
    },
  });

  async function cookPlease(data: SignupData) {
    console.log(data);
    setLoading(true);
    const response = await signup(data);
    setError(response.data?.message ?? response.error ?? "");
    setLoading(false);
  }

  const formId = useId();

  return (
    <Form {...form}>
      <Card className="mx-auto md:my-3 lg:w-1/2 md:w-3/4">
        <CardHeader>
          <CardTitle>Create an Account</CardTitle>
        </CardHeader>
        <CardContent asChild>
          <form
            id={formId}
            onSubmit={form.handleSubmit(cookPlease)}
            className="gap-y-2 space-y-3"
            {...props}
          >
            <SignupField name="firstName" label="First Name" />
            <SignupField name="lastName" label="Last Name" />
            <SignupField name="email" label="Email" />
            <SignupField type="password" name="password" label="Password" />
            <SignupField type="password" name="passwordConfirm" label="Confirm Password" />
          </form>
        </CardContent>
        <CardFooter>
          <Button disabled={loading} form={formId} type="submit">
            Create
          </Button>
          <p role="alert">{error}</p>
        </CardFooter>
      </Card>
    </Form>
  );
}
