import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Form, FormField, FormMessage, FormSubmit } from "@/components/ui/form-server";

const ResetForm = Form<ResetLinkData>;
const ResetField = FormField<ResetLinkData>;
const ResetMessage = FormMessage<ResetLinkData>;

import { Input } from "@/components/ui/input";
import { ResetLinkData } from "@/lib/schema/auth";
import { sendResetLink } from "@/lib/actions/auth";

export default function Reset() {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <ResetForm action={sendResetLink} className="flex flex-col">
          <CardHeader>
            <CardTitle>Reset Password</CardTitle>
            <CardDescription>
              Enter your email, and we will send you a password reset link for your account
            </CardDescription>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-x-3 gap-y-4">
            <ResetField name="email" label="Email">
              <Input autoComplete="email" id="email" name="email" placeholder="Email" />
            </ResetField>
          </CardContent>
          <CardFooter className="mt-auto flex justify-between">
            <ResetMessage />
            <FormSubmit>Send reset link</FormSubmit>
          </CardFooter>
        </ResetForm>
      </Card>
    </main>
  );
}
