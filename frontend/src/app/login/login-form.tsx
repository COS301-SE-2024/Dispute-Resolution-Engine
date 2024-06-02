"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import React, { HTMLAttributes, useId, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import { login } from "@/app/lib/auth";
import TextField from "@/components/form/text-field";
import { useRouter } from "next/navigation";

const loginSchema = z.object({
  email: z.string().min(1, "Required").email("Please enter a valid email"),
  password: z.string().min(1, "Required"),
});

export type LoginData = z.infer<typeof loginSchema>;
const LoginField = TextField<LoginData>;

export default function LoginForm(props: HTMLAttributes<HTMLFormElement>) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const router = useRouter();

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
    const res = await login(data);
    setLoading(false);
    if (res.error) {
      setError(res.error);
      return false;
    }
    router.push("/disputes");
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
            <LoginField name="email" label="Email" />
            <LoginField type="password" name="password" label="Password" />
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
