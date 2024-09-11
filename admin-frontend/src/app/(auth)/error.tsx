"use client";

import { Button } from "@/components/ui/button";

import Link from "next/link";

export default function Error({ error }: { error: Error & { digest?: string } }) {
  return (
      <div className="flex flex-col items-center justify-center mt-10">
        <main>
          <h2 className="font-bold text-3xl mb-2">Something went wrong.</h2>
          <p className="mb-4">{error.message}</p>
          <Button asChild variant="outline">
              <Link href="/login">
                Back to Login
              </Link>
          </Button>
        </main>
      </div>
  );
}
