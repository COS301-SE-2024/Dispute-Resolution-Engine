"use client";

import { signupSchema, type SignupData } from "@/lib/schema/auth";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronRight, User } from "lucide-react";
import { HTMLAttributes, ReactNode, useId, useState } from "react";
import { Controller, FieldValues, FormProvider, useForm, useFormContext } from "react-hook-form";

import * as RadioGroup from "@radix-ui/react-radio-group";
import { Input } from "../ui/input";
import GenderSelect from "./gender-select";
import LanguageSelect from "./language.select";
import CountrySelect from "./country-select";
import { Label } from "../ui/label";
import { FormField, FormMessage } from "../ui/form-client";
import { Button } from "../ui/button";
import { signup } from "@/lib/actions/auth";

const steps: {
  id: string;
  name: string;
  fields: (keyof SignupData)[];
}[] = [
  { id: "Step 1", name: "User Type", fields: ["userType"] },
  { id: "Step 2", name: "The Basics", fields: ["email", "password", "passwordConfirm"] },
  {
    id: "Step 3",
    name: "Personal Details",
    fields: ["firstName", "lastName", "gender", "nationality", "preferredLanguage", "dateOfBirth"],
  },
];

const SignupMessage = FormMessage<SignupData>;
const SignupField = FormField<SignupData>;

export default function SignupForm() {
  const emailId = useId();
  const passId = useId();
  const confirmId = useId();

  const fnameId = useId();
  const lnameId = useId();
  const dobId = useId();
  const genderId = useId();
  const countryId = useId();
  const langId = useId();

  const form = useForm<SignupData>({
    resolver: zodResolver(signupSchema),
  });
  const { setError, register, handleSubmit, control, trigger } = form;

  const [currentStep, setCurrentStep] = useState(0);

  async function onSubmit(form: SignupData) {
    const res = await signup(form);
    if (res.error) {
      setError("root", { type: "custom", message: res.error });
    }
  }

  async function nav(index: number) {
    if (index == currentStep) {
      return;
    } else if (index > currentStep) {
      for (let i = currentStep; i < index; i++) {
        if (!(await trigger(steps[currentStep].fields, { shouldFocus: true }))) {
          setCurrentStep(i);
          return;
        }
      }
      setCurrentStep(index);
    } else {
      trigger(steps[currentStep].fields, { shouldFocus: true });
      setCurrentStep(index);
    }
  }

  return (
    <FormProvider {...form}>
      <div className="mx-auto w-fit">
        <nav>
          <ol className="flex gap-3">
            {steps.map((step, i) => (
              <li key={step.id}>
                <SignupStep
                  onClick={() => nav(i)}
                  name={step.id}
                  desc={step.name}
                  active={i <= currentStep}
                />
              </li>
            ))}
          </ol>
        </nav>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {currentStep == 0 && (
            <>
              <Controller
                name="userType"
                control={control}
                rules={{ required: true }}
                render={({ field }) => {
                  const { onChange, ...field2 } = field;
                  return (
                    <RadioGroup.Root
                      onValueChange={onChange}
                      {...field2}
                      className="flex flex-col gap-4"
                    >
                      <RadioGroup.Item
                        value="user"
                        asChild
                        className="data-[state='checked']:border-dre-100"
                      >
                        <Button variant="outline">
                          <h2>User</h2>
                        </Button>
                      </RadioGroup.Item>
                      <RadioGroup.Item
                        value="expert"
                        asChild
                        className="data-[state='checked']:border-dre-100"
                      >
                        <Button variant="outline">
                          <h2>Expert</h2>
                        </Button>
                      </RadioGroup.Item>
                    </RadioGroup.Root>
                  );
                }}
              />
              <div className="flex justify-end">
                <Button
                  type="button"
                  aria-label="Next"
                  title="Next"
                  variant="outline"
                  className="ml-auto"
                  onClick={() => nav(1)}
                >
                  <ChevronRight />
                </Button>
              </div>
            </>
          )}
          {currentStep == 1 && (
            <>
              <SignupField id={fnameId} name="firstName" label="First Name">
                <Input
                  id={fnameId}
                  {...register("firstName")}
                  autoComplete="given-name"
                  placeholder="First Name"
                />
              </SignupField>
              <SignupField id={lnameId} name="lastName" label="Last Name">
                <Input
                  id={lnameId}
                  {...register("lastName")}
                  autoComplete="family-name"
                  placeholder="Last Name"
                />
              </SignupField>
              <SignupField id={emailId} name="email" label="Email">
                <Input
                  id={emailId}
                  autoComplete="email"
                  placeholder="Email"
                  {...register("email")}
                />
              </SignupField>
              <SignupField id={passId} name="password" label="Password">
                <Input
                  id={passId}
                  autoComplete="new-password"
                  placeholder="Password"
                  type="password"
                  {...register("password")}
                />
              </SignupField>
              <SignupField id={confirmId} name="passwordConfirm" label="Confirm Password">
                <Input
                  id={confirmId}
                  autoComplete="new-password"
                  placeholder="Confirm Password"
                  type="password"
                  {...register("passwordConfirm")}
                />
              </SignupField>
              <div className="flex justify-end">
                <Button
                  type="button"
                  aria-label="Next"
                  title="Next"
                  variant="outline"
                  className="ml-auto"
                  onClick={() => nav(2)}
                >
                  <ChevronRight />
                </Button>
              </div>
            </>
          )}
          {currentStep == 2 && (
            <>
              <SignupField id={genderId} name="gender" label="Gender">
                <Controller
                  name="gender"
                  control={control}
                  rules={{ required: true }}
                  render={({ field }) => {
                    const { onChange, ...field2 } = field;
                    return <GenderSelect id={genderId} onValueChange={onChange} {...field2} />;
                  }}
                />
              </SignupField>
              <SignupField id={langId} name="preferredLanguage" label="Preferred Language">
                <Controller
                  name="preferredLanguage"
                  control={control}
                  render={({ field }) => {
                    const { onChange, ...field2 } = field;
                    return <LanguageSelect id={langId} onValueChange={onChange} {...field2} />;
                  }}
                />
              </SignupField>
              <SignupField id={countryId} name="nationality" label="Nationality">
                <Controller
                  name="nationality"
                  control={control}
                  rules={{ required: true }}
                  render={({ field }) => {
                    const { onChange, ...field2 } = field;
                    return <CountrySelect id={countryId} onValueChange={onChange} {...field2} />;
                  }}
                />
              </SignupField>
              <SignupField id={dobId} name="dateOfBirth" label="Date of Birth">
                <Input
                  id={dobId}
                  {...register("dateOfBirth")}
                  autoComplete="bday"
                  type="date"
                  className="w-fit"
                />
              </SignupField>
              <div className="col-span-2 flex justify-end items-center gap-2">
                <SignupMessage />
                <Button type="submit">Sign Up</Button>
              </div>
            </>
          )}
        </form>
      </div>
    </FormProvider>
  );
}

function SignupStep({
  name,
  desc,
  active,
  onClick = () => {},
}: {
  name: string;
  desc: string;
  active: boolean;
  onClick?: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "text-left py-2 border-t-4 w-44",
        active ? " border-dre-200" : "border-dre-bg-light/50",
      )}
    >
      <h3>{name}</h3>
      <p>{desc}</p>
    </button>
  );
}

function Navbar({
  page,
  maxPages,
  onNavigate,
}: {
  page: number;
  maxPages: number;
  onNavigate: (step: number) => void;
}) {
  return (
    <div className="col-span-2 flex justify-between">
      <Button type="button" disabled={page == 0} aria-label="Back" title="Back">
        Back
      </Button>
      <Button type="button" disabled={page == maxPages - 1} aria-label="Next" title="Next">
        Next
      </Button>
    </div>
  );
}
