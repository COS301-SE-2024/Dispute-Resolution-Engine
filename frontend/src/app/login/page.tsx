import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import { Form, FormField, FormMessage, FormSubmit } from "@/components/ui/form-server";
import { LoginData } from "@/lib/schema/auth";
import { login } from "@/lib/actions/auth";
import { useId } from "react";
import { cn } from "@/lib/utils";

const LoginForm = Form<LoginData>;
const LoginField = FormField<LoginData>;
const LoginMessage = FormMessage<LoginData>;

export default function Login() {
  const emailId = useId();
  const passId = useId();

  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <LoginForm className="flex flex-col" action={login}>
          <CardHeader>
            <CardTitle>Login</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <LoginField id={emailId} name="email" label="Email">
              <Input id={emailId} name="email" autoComplete="email" placeholder="Email" />
            </LoginField>
            <LoginField id={passId} name="password" label="Password">
              <Input
                id={passId}
                name="password"
                autoComplete="current-password"
                placeholder="Password"
                type="password"
              />
              <Link href="/reset" className={buttonVariants({ variant: "link" })}>
                Forgot Password?
              </Link>
            </LoginField>
          </CardContent>
          <CardFooter className="mt-auto flex justify-between">
            <p>
              {"Don't have an account? "}
              <Link href="/signup" className="hover:underline dark:text-dre-100 text-dre-200">
                Create one
              </Link>
            </p>
            <LoginMessage />
            <FormSubmit>Login</FormSubmit>
          </CardFooter>
        </LoginForm>
      </Card>
    </main>
  );
}
