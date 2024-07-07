import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { SignupButton, SignupField, SignupForm, SignupMessage } from "./signup-form";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import CountrySelect from "@/components/form/country-select";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import { SignupError } from "@/lib/schema/auth";
import { Result } from "@/lib/interfaces";
import { signup } from "@/lib/actions/auth";

export default function Signup() {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <SignupForm className="flex flex-col">
          <CardHeader>
            <CardTitle>Signup</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-x-3 gap-y-4">
            <SignupField name="firstName" label="First Name">
              <Input
                id="firstName"
                name="firstName"
                autoComplete="given-name"
                placeholder="First Name"
              />
            </SignupField>
            <SignupField name="lastName" label="Last Name">
              <Input
                id="lastName"
                name="lastName"
                autoComplete="family-name"
                placeholder="Last Name"
              />
            </SignupField>
            <SignupField name="email" label="Email">
              <Input autoComplete="email" id="email" name="email" placeholder="Email" />
            </SignupField>
            <SignupField name="dateOfBirth" label="Date of Birth">
              <Input
                id="dateOfBirth"
                name="dateOfBirth"
                autoComplete="bday"
                type="date"
                className="w-fit"
              />
            </SignupField>
            <SignupField name="password" label="Password" className="col-span-2">
              <Input
                autoComplete="new-password"
                id="password"
                name="password"
                placeholder="Password"
                type="password"
              />
            </SignupField>
            <SignupField name="passwordConfirm" label="Confirm Password" className="col-span-2">
              <Input
                autoComplete="new-password"
                id="passwordConfirm"
                name="passwordConfirm"
                placeholder="Confirm Password"
                type="password"
              />
            </SignupField>
          </CardContent>
          <CardFooter className="mt-auto flex justify-between">
            <p>
              {"Already have a account?"}
              <Link href="/login" className={buttonVariants({ variant: "link" })}>
                Login
              </Link>
            </p>
            <SignupMessage />
            <SignupButton />
          </CardFooter>
        </SignupForm>
      </Card>
    </main>
  );
}
