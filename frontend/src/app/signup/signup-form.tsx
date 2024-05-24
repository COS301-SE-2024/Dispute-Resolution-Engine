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
import { ReactElement, JSXElementConstructor } from "react";
import {
  ControllerFieldState,
  ControllerRenderProps,
  FieldValues,
  UseFormStateReturn,
  useForm,
} from "react-hook-form";
import { z } from "zod";

const signupSchema = z.object({
  email: z.string().min(1, "Required").email("Please enter a valid email"),
  password: z.string().min(1, "Required"),
  passConfirm: z.string().min(1, "Required"),
});

type SignupData = z.infer<typeof signupSchema>;

export default function SignupForm() {
  const form = useForm<SignupData>({
    resolver: zodResolver(signupSchema),
    defaultValues: {
      email: "",
      password: "",
      passConfirm: "",
    },
  });

  function onSubmit(values: SignupData) {
    console.log(values);
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="mx-auto my-20 w-fit flex flex-col gap-2"
      >
        <FormField
          name="email"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Email</FormLabel>
              <FormControl>
                <input
                  placeholder="Email"
                  {...field}
                  className="py-1 px-2 border-[1px] border-gray-200 rounded-md container"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="password"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Password</FormLabel>
              <FormControl>
                <input
                  type="password"
                  placeholder="Password"
                  {...field}
                  className="py-1 px-2 border-[1px] border-gray-200 rounded-md container"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          name="passConfirm"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <input
                  type="password"
                  placeholder="Confirm Password"
                  {...field}
                  className="py-1 px-2 border-[1px] border-gray-200 rounded-md container"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Signup</Button>
      </form>
    </Form>
  );
}
