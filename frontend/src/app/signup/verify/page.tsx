import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from "@/components/ui/input-otp";
import Link from "next/link";
import { useId } from "react";
import { VerifyButton, VerifyField, VerifyForm, VerifyMessage } from "./verify-form";

export default function Verify() {
  const formId = useId();

  return (
    <main className="flex flex-col justify-center items-center h-full gap-5">
      <div className="text-center">
        <CardTitle>Check your email</CardTitle>
        <CardDescription>We sent an OTP to your email address</CardDescription>
      </div>
      <VerifyForm className="flex flex-col justify-center items-center gap-3">
        <VerifyField name="pin" label="Pin" className="flex flex-col items-center">
          <InputOTP maxLength={6} name="pin">
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
        <VerifyButton />
        <VerifyMessage />
      </VerifyForm>
    </main>
  );
}
