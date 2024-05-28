"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import React, { HTMLAttributes, useId, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { signup } from "../lib/signup";
import { useFormState } from "react-dom";

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

function TextField({ name, label }: { name: keyof SignupData; label: string }) {
  return (
    <FormField
      name={name}
      render={({ field }) => (
        <FormItem>
          <FormLabel className="block">{label}</FormLabel>
          <FormControl>
            <input
              placeholder={label}
              {...field}
              className="w-full py-2 px-3 border-[1px] border-gray-200 rounded-md"
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

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
    setError(await signup(data));
    setLoading(false);
  }

  const formId = useId();

  return (
    <Form {...form}>
      <Card className="mx-auto md:my-3 lg:w-1/2 md:w-3/4">
        <CardHeader>
          <CardTitle>Create an Account</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            id={formId}
            onSubmit={form.handleSubmit(cookPlease)}
            className="gap-y-2 space-y-3"
            {...props}
          >
            <TextField name="firstName" label="First Name" />
            <TextField name="lastName" label="Last Name" />
            <TextField name="email" label="Email" />
            <TextField name="password" label="Password" />
            <TextField name="passwordConfirm" label="Confirm Password" />
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
