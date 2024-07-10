import { Form, FormField, FormMessage, FormSubmit } from "@/components/ui/form-server";
import { CardDescription, CardTitle } from "@/components/ui/card";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from "@/components/ui/input-otp";

import { verify } from "@/lib/actions/auth";
import { VerifyData } from "@/lib/schema/auth";
import { useId } from "react";

const VerifyForm = Form<VerifyData>;
const VerifyMessage = FormMessage<VerifyData>;
const VerifyField = FormField<VerifyData>;

export default function Verify() {
  const pinId = useId();

  return (
    <main className="flex flex-col justify-center items-center h-full gap-5">
      <div className="text-center">
        <CardTitle>Check your email</CardTitle>
        <CardDescription>We sent an OTP to your email address</CardDescription>
      </div>
      <VerifyForm action={verify} className="flex flex-col justify-center items-center gap-3">
        <VerifyField id={pinId} name="pin" label="Pin" className="flex flex-col items-center">
          <InputOTP maxLength={6} name="pin" id={pinId}>
            <InputOTPGroup>
              <InputOTPSlot index={0} />
              <InputOTPSlot index={1} />
              <InputOTPSlot index={2} />
            </InputOTPGroup>
            <InputOTPSeparator />
            <InputOTPGroup>
              <InputOTPSlot index={3} />
              <InputOTPSlot index={4} />
              <InputOTPSlot index={5} />
            </InputOTPGroup>
          </InputOTP>
        </VerifyField>
        <FormSubmit>Verify</FormSubmit>
        <VerifyMessage />
      </VerifyForm>
    </main>
  );
}
