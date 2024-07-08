import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";
import { Form, FormField, FormMessage, FormSubmit } from "@/components/form/form";
import { LoginData } from "@/lib/schema/auth";
import { login } from "@/lib/actions/auth";

const LoginForm = Form<LoginData>;
const LoginField = FormField<LoginData>;
const LoginMessage = FormMessage<LoginData>;

export default function Login() {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <LoginForm className="flex flex-col" action={login}>
          <CardHeader>
            <CardTitle>Login</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <LoginField name="email" label="Email">
              <Input autoComplete="email" id="email" name="email" placeholder="Email" />
            </LoginField>
            <LoginField name="password" label="Password">
              <Input
                id="password"
                name="password"
                autoComplete="current-password"
                placeholder="Password"
                type="password"
              />
            </LoginField>
          </CardContent>
          <CardFooter className="mt-auto flex justify-between">
            <p>
              {"Don't have an account?"}
              <Link href="/signup" className={buttonVariants({ variant: "link" })}>
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
