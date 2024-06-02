import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
  InputOTPSeparator,
} from "@/components/ui/input-otp";
import Link from "next/link";

export default function Verify() {
  return (
    <main>
      <Card className="mx-auto md:my-3 lg:w-1/2 md:w-3/4">
        <CardHeader>
          <CardTitle>Check your email</CardTitle>
          <CardDescription>We sent an OTP to your email address</CardDescription>
        </CardHeader>
        <CardContent className="mx-auto">
          <InputOTP maxLength={6}>
            <InputOTPGroup>
              <InputOTPSlot index={0} />
              <InputOTPSlot index={1} />
              <InputOTPSlot index={2} />
            </InputOTPGroup>
            <InputOTPSeparator />
            <InputOTPGroup>
              <InputOTPSlot index={3} />
              <InputOTPSlot index={4} />
              <InputOTPSlot index={5} />
            </InputOTPGroup>
          </InputOTP>
        </CardContent>
        <CardFooter>
          <Button asChild>
            <Link href="/login">Submit</Link>
          </Button>
        </CardFooter>
      </Card>
    </main>
  );
}
