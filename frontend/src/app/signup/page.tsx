"use client";

import { Button } from "@/components/ui/button";
import { CardDescription, CardTitle } from "@/components/ui/card";
import { Providers, SignupForm, SignupMessage } from "./signup-form";

import Link from "next/link";
import { useId } from "react";

export default function Signup() {
  const formId = useId();

  return (
    <Providers>
      <div className="h-full overflow-y-auto">
        <div className="mx-auto md:w-[40rem] grid grid-rows-[auto_1fr_auto] gap-5 px-5  h-full">
          <header className="space-y-2 mt-5">
            <CardTitle>Signup</CardTitle>
            <CardDescription>Welcome to DRE! Tell us a little bit about yourself</CardDescription>
          </header>
          <SignupForm id={formId} />
          <footer className="flex items-center flex-wrap gap-2 pb-4">
            <p className="grow">
              Already have an account?{" "}
              <Link href="/login" className="hover:underline dark:text-dre-100 text-dre-200">
                Login
              </Link>
            </p>
            <SignupMessage />
            <Button type="submit" form={formId}>
              Signup
            </Button>
          </footer>
        </div>
      </div>
    </Providers>
  );
}
