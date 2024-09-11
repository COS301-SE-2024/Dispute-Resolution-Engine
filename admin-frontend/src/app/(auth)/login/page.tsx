import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import Link from "next/link";
import { useId } from "react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { Button } from "@/components/ui/button";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Login",
  icons: "/logo.svg",
};

export default function Login() {
  const emailId = useId();
  const passId = useId();

  return (
  <Card className="md:mx-auto md:max-w-xl mt-5 mx-2 sm:mx-5">
        <form>
          <CardHeader>
            <CardTitle>Admin Login</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
          <div>
          <Label htmlFor={emailId}>Email</Label>
          <Input id={emailId} name="email" autoComplete="email" placeholder="Email" />
          </div>

          <div>
          <Label htmlFor={passId}>Password</Label>
          <Input
          id={passId}
          name="password"
          autoComplete="current-password"
          placeholder="Password"
          type="password"
          />
          </div>
          </CardContent>
          <CardFooter className="mt-auto flex justify-end">
            <Button className="ml-auto">Login</Button>
          </CardFooter>
        </form>
  </Card>
  );
}
