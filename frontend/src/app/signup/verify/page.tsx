import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { useId } from "react";
import { Button } from "@/components/ui/button";
import { Providers, VerifyForm, VerifyMessage } from "./verify-form";
import ResendForm from "./resend-form";

export default function Verify() {
  const formId = useId();

  return (
    <Providers>
      <div className="flex flex-col justify-center items-center h-full gap-5">
        <Card className="mx-auto md:w-[30rem]">
          <CardHeader>
            <CardTitle>Confirm your email</CardTitle>
            <CardDescription>Check your email! We sent you a code.</CardDescription>
          </CardHeader>
          <CardContent className="flex justify-center">
            <VerifyForm id={formId} />
            <VerifyMessage />
          </CardContent>
          <CardFooter asChild className="flex items-center flex-wrap gap-2">
            <footer>
              <p className="grow">
                {"Didn't receive the code? "}
                <ResendForm />
              </p>
              <Button form={formId} type="submit">
                Confirm
              </Button>
            </footer>
          </CardFooter>
        </Card>
      </div>
    </Providers>
  );
}
