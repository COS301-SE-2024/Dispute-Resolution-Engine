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
import { login, signup } from "@/app/lib/auth";
import { Input } from "@/components/ui/input";

const loginSchema = z.object({
  email: z.string().min(1, "Required").email("Please enter a valid email"),
  password: z.string().min(1, "Required"),
});

export type LoginData = z.infer<typeof loginSchema>;

function TextField({ name, label }: { name: keyof LoginData; label: string }) {
  return (
    <FormField
      name={name}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{label}</FormLabel>
          <FormControl>
            <Input placeholder={label} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

export default function LoginForm(props: HTMLAttributes<HTMLFormElement>) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const form = useForm<LoginData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function cookPlease(data: LoginData) {
    console.log(data);
    setLoading(true);
    setError(await login(data));
    setLoading(false);
  }

  const formId = useId();

  return (
    <Form {...form}>
      <Card className="mx-auto md:my-3 lg:w-1/2 md:w-3/4">
        <CardHeader>
          <CardTitle>Login</CardTitle>
        </CardHeader>
        <CardContent asChild>
          <form
            id={formId}
            onSubmit={form.handleSubmit(cookPlease)}
            className="gap-y-2 space-y-3"
            {...props}
          >
            <TextField name="email" label="Email" />
            <TextField name="password" label="Password" />
          </form>
        </CardContent>
        <CardFooter>
          <Button disabled={loading} form={formId} type="submit">
            Login
          </Button>
          <p role="alert">{error}</p>
        </CardFooter>
      </Card>
    </Form>
  );
}
