"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useId, type ReactNode } from "react";
import { Controller, FormProvider, useForm, useFormContext } from "react-hook-form";

import { verifySchema, type VerifyData } from "@/lib/schema/auth";
import { FormField, FormMessage } from "@/components/ui/form-client";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "@/components/ui/input-otp";
import { verify } from "@/lib/actions/auth";

export const VerifyMessage = FormMessage<VerifyData>;

export function Providers({ children }: { children?: ReactNode | ReactNode[] }) {
  const form = useForm<VerifyData>({
    resolver: zodResolver(verifySchema),
  });
  return <FormProvider {...form}>{children}</FormProvider>;
}

export function VerifyForm({ id }: { id?: string }) {
  const { control, setError, handleSubmit } = useFormContext<VerifyData>();
  async function onSubmit(data: VerifyData) {
    console.log("submit");
    await verify(data).catch((e: Error) => {
      setError("root", {
        type: "custom",
        message: e.message,
      });
    });
  }

  const pinId = useId();

  return (
    <form
      id={id}
      className="flex flex-col justify-center items-center gap-3"
      onSubmit={handleSubmit(onSubmit)}
    >
      <Controller
        name="pin"
        control={control}
        rules={{ required: true }}
        render={({ field }) => (
          <InputOTP maxLength={6} id={pinId} {...field}>
            <InputOTPGroup>
              <InputOTPSlot index={0} />
              <InputOTPSlot index={1} />
              <InputOTPSlot index={2} />
              <InputOTPSlot index={3} />
              <InputOTPSlot index={4} />
              <InputOTPSlot index={5} />
            </InputOTPGroup>
          </InputOTP>
        )}
      />

      <VerifyMessage name="pin" />
    </form>
  );
}
