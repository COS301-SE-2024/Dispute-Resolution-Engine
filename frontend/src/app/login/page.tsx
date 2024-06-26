import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { LoginForm, LoginField, LoginMessage, LoginButton } from "./login-form";
import Link from "next/link";
import { buttonVariants } from "@/components/ui/button";

export default function Login() {
  return (
    <main className="md:pt-3 h-full">
      <Card variant="page" asChild>
        <LoginForm className="flex flex-col">
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
            <LoginButton />
          </CardFooter>
        </LoginForm>
      </Card>
    </main>
  );
}
