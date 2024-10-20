"use client";

import CountrySelect from "@/components/form/country-select";
import { FormField, FormMessage } from "@/components/ui/form-client";
import { Input } from "@/components/ui/input";
import { signup } from "@/lib/actions/auth";
import { SignupData, signupSchema } from "@/lib/schema/auth";

import { zodResolver } from "@hookform/resolvers/zod";
import { ReactNode, useId } from "react";
import { Controller, FormProvider, useForm, useFormContext } from "react-hook-form";

export const SignupMessage = FormMessage<SignupData>;
export const SignupField = FormField<SignupData>;

export function Providers({ children }: { children?: ReactNode | ReactNode[] }) {
  const form = useForm<SignupData>({
    resolver: zodResolver(signupSchema),
  });
  return <FormProvider {...form}>{children}</FormProvider>;
}

export function SignupForm({ id }: { id?: string }) {
  const { control, register, handleSubmit, setError } = useFormContext<SignupData>();

  async function onSubmit(data: SignupData) {
    const { error } = await signup(data);
    if (error) {
      setError("root", {
        type: "custom",
        message: error,
      });
    }
  }

  const emailId = useId();
  const passId = useId();
  const confirmId = useId();

  const fnameId = useId();
  const lnameId = useId();
  const countryId = useId();

  return (
    <form className="space-y-5" id={id} onSubmit={handleSubmit(onSubmit)}>
      <div className="grid sm:grid-cols-2 gap-x-2 gap-y-5">
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
      </div>
      <SignupField id={emailId} name="email" label="Email">
        <Input id={emailId} autoComplete="email" placeholder="Email" {...register("email")} />
      </SignupField>

      <SignupField id={countryId} name="nationality" label="Nationality">
        <Controller
          name="nationality"
          control={control}
          rules={{ required: true }}
          render={({ field }) => {
            const { onChange, ref, ...field2 } = field;
            return <CountrySelect id={countryId} onValueChange={onChange} {...field2} />;
          }}
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
    </form>
  );
}
