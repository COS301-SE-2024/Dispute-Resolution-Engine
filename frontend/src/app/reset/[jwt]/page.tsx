import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
  CardDescription,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Form, FormField, FormMessage, FormSubmit } from "@/components/ui/form-server";
import { ResetPassData } from "@/lib/schema/auth";
import { resetPassword } from "@/lib/actions/auth";
import { useId } from "react";

const ResetForm = Form<ResetPassData>;
const ResetField = FormField<ResetPassData>;
const ResetMessage = FormMessage<ResetPassData>;

type Props = {
  params: { jwt: string };
};

export default function ResetPage({ params }: Props) {
  const newId = useId();
  const confirmId = useId();

  return (
    <main className="md:pt-3 h-full">
      <Card asChild>
        <ResetForm action={resetPassword} className="flex flex-col">
          <CardHeader>
            <CardTitle>Reset Password</CardTitle>
            <CardDescription>Enter your new password</CardDescription>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-x-3 gap-y-4">
            <input type="hidden" name="jwt" value={params.jwt} />
            <ResetField id={newId} name="password" label="Password" className="col-span-2">
              <Input
                id={newId}
                name="password"
                autoComplete="new-password"
                placeholder="Password"
                type="password"
              />
            </ResetField>
            <ResetField
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
            </ResetField>
          </CardContent>
          <CardFooter className="mt-auto flex justify-between">
            <ResetMessage />
            <FormSubmit>Reset Password</FormSubmit>
          </CardFooter>
        </ResetForm>
      </Card>
    </main>
  );
}
