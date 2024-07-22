"use client";

import { FormSubmit } from "@/components/ui/form-server";
import { resendOTP } from "@/lib/actions/auth";
import { useFormState } from "react-dom";

export default function ResendForm() {
  const [_state, formAction] = useFormState(resendOTP, undefined);

  return (
    <form action={formAction}>
      <FormSubmit variant="link">Resend OTP</FormSubmit>
    </form>
  );
}
