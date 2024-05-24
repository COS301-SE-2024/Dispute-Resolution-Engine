"use client";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const signupSchema = z
  .object({
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

type SignupData = z.infer<typeof signupSchema>;

export default function SignupForm() {
  const form = useForm<SignupData>({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      email: "",
      password: "",
      passwordConfirm: "",
    },
  });

  function onSubmit(values: SignupData) {
    console.log(values);
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="bg-white rounded-xl mx-auto my-10 w-1/2 p-5 shadow-lg gap-2 space-y-3"
      >
        <h1 className="text-lg font-medium mb-3">Create an Account</h1>
        <FormField
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="block">Email</FormLabel>
              <FormControl>
                <input
                  placeholder="Email"
                  {...field}
                  className="w-full py-2 px-3 border-[1px] border-gray-200 rounded-md"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <input
                  type="password"
                  placeholder="Password"
                  {...field}
                  className="w-full py-2 px-3 border-[1px] border-gray-200 rounded-md"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="passwordConfirm"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <input
                  type="password"
                  placeholder="Confirm Password"
                  {...field}
                  className="w-full py-2 px-3 border-[1px] border-gray-200 rounded-md"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit">Create</Button>
      </form>
    </Form>
  );
}
