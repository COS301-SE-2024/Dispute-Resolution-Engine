import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { LoginForm, LoginField, LoginMessage, LoginButton } from "./login-form";

export default function Login() {
  return (
    <main className="md:pt-3 h-full">
      <Card asChild>
        <LoginForm className="mx-auto lg:w-1/2 md:w-3/4 md:h-fit h-full flex flex-col">
          <CardHeader>
            <CardTitle>Login</CardTitle>
          </CardHeader>
          <CardContent>
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
          <CardFooter className="mt-auto flex md:justify-start justify-end">
            <LoginMessage />
            <LoginButton />
          </CardFooter>
        </LoginForm>
      </Card>
    </main>
  );
}
