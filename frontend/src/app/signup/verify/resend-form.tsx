"use client";

import { resendOTP } from "@/lib/actions/auth";
import { useFormState } from "react-dom";

export default function ResendForm() {
  const [_state, formAction] = useFormState(resendOTP, undefined);

  return (
    <form action={formAction} className="inline">
      <button className="text-dre-200 dark:text-dre-100 hover:underline">Resend code</button>
    </form>
  );
}
