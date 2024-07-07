import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ResetForm, ResetButton, ResetField, ResetMessage } from "./reset-form";
import { Input } from "@/components/ui/input";

export default function Reset() {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <ResetForm className="flex flex-col">
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
            <ResetButton />
          </CardFooter>
        </ResetForm>
      </Card>
    </main>
  );
}
