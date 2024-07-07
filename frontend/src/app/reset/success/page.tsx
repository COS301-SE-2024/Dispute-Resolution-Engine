import { CardDescription, CardTitle } from "@/components/ui/card";

export default function ResetSent() {
  return (
    <main className="flex flex-col justify-center items-center h-full gap-5">
      <div className="text-center">
        <CardTitle>Check your email</CardTitle>
        <CardDescription>We sent a password reset link to your email address</CardDescription>
      </div>
    </main>
  );
}
