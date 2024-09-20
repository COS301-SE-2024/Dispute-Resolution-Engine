"use client";

import { Button } from "@/components/ui/button";
import Navbar from "@/components/navbar";

import { useRouter } from "next/navigation";

export default function Error({ error }: { error: Error & { digest?: string } }) {
  const router = useRouter();

  return (
    <div className="flex flex-col items-center justify-center h-full">
      <main>
        <h2 className="font-bold text-3xl mb-2">Something went wrong.</h2>
        <p className="mb-4">{error.message}</p>
        <Button onClick={() => router.back()} variant="outline">
          Back
        </Button>
      </main>
    </div>
  );
}
