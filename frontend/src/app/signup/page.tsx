import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import CountrySelect from "@/components/form/country-select";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import GenderSelect from "@/components/form/gender-select";
import LanguageSelect from "@/components/form/language.select";
import { Form, FormField, FormMessage, FormSubmit } from "@/components/ui/form-server";
import { SignupData } from "@/lib/schema/auth";
import { signup } from "@/lib/actions/auth";
import { useId } from "react";

const SignupForm = Form<SignupData>;
const SignupMessage = FormMessage<SignupData>;
const SignupField = FormField<SignupData>;

export default function Signup() {
  const fnameId = useId();
  const lnameId = useId();
  const emailId = useId();
  const dobId = useId();
  const genderId = useId();
  const countryId = useId();
  const langId = useId();
  const passId = useId();
  const confirmId = useId();

  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <SignupForm action={signup} className="flex flex-col">
          <CardHeader>
            <CardTitle>Signup</CardTitle>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-x-3 gap-y-4">
            <SignupField id={fnameId} name="firstName" label="First Name">
              <Input
                id={fnameId}
                name="firstName"
                autoComplete="given-name"
                placeholder="First Name"
              />
            </SignupField>
            <SignupField id={lnameId} name="lastName" label="Last Name">
              <Input
                id={lnameId}
                name="lastName"
                autoComplete="family-name"
                placeholder="Last Name"
              />
            </SignupField>
            <SignupField id={emailId} name="email" label="Email">
              <Input id={emailId} name="email" autoComplete="email" placeholder="Email" />
            </SignupField>
            <SignupField id={dobId} name="dateOfBirth" label="Date of Birth">
              <Input
                id={dobId}
                name="dateOfBirth"
                autoComplete="bday"
                type="date"
                className="w-fit"
              />
            </SignupField>
            <SignupField id={genderId} name="gender" label="Gender">
              <GenderSelect id={genderId} name="gender" />
            </SignupField>
            <SignupField id={langId} name="preferredLanguage" label="Preferred Language">
              <LanguageSelect id={langId} name="preferredLanguage" />
            </SignupField>
            <SignupField id={countryId} name="nationality" label="Nationality">
              <CountrySelect id={countryId} name="nationality" />
            </SignupField>
            <SignupField id={passId} name="password" label="Password" className="col-span-2">
              <Input
                id={passId}
                name="password"
                autoComplete="new-password"
                placeholder="Password"
                type="password"
              />
            </SignupField>
            <SignupField
              id={confirmId}
              name="passwordConfirm"
              label="Confirm Password"
              className="col-span-2"
            >
              <Input
                id={confirmId}
                name="passwordConfirm"
                autoComplete="new-password"
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
            <FormSubmit>Signup</FormSubmit>
          </CardFooter>
        </SignupForm>
      </Card>
    </main>
  );
}
