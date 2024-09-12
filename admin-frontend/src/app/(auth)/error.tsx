"use client";

import { Button } from "@/components/ui/button";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="flex flex-col items-center justify-center mt-10">
      <main>
        <h2 className="font-bold text-3xl mb-2">Something went wrong.</h2>
        <p className="mb-4">{error.message}</p>
        <Button onClick={() => reset()} variant="outline">
          Reload
        </Button>
      </main>
    </div>
  );
}
