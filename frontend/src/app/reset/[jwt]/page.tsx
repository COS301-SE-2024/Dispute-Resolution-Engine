import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
  CardDescription,
} from "@/components/ui/card";
import { ResetForm, ResetField, ResetMessage, ResetButton } from "./reset-form";
import { Input } from "@/components/ui/input";

type Props = {
  params: { jwt: string };
};

export default function ResetPage({ params }: Props) {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <ResetForm className="flex flex-col">
          <CardHeader>
            <CardTitle>Reset Password</CardTitle>
            <CardDescription>Enter your new password</CardDescription>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-x-3 gap-y-4">
            <input type="hidden" name="jwt" value={params.jwt} />
            <ResetField name="password" label="Password" className="col-span-2">
              <Input
                autoComplete="new-password"
                id="password"
                name="password"
                placeholder="Password"
                type="password"
              />
            </ResetField>
            <ResetField name="passwordConfirm" label="Confirm Password" className="col-span-2">
              <Input
                autoComplete="new-password"
                id="passwordConfirm"
                name="passwordConfirm"
                placeholder="Confirm Password"
                type="password"
              />
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
