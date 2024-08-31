"use client";
import { Button } from "@/components/ui/button";
import { useEffect } from "react";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Log the error to an error reporting service
    console.error(error);
  }, [error]);

  return (
    <div className="flex flex-col items-center justify-center">
      <main>
        <h2 className="font-bold text-3xl mb-2">500 Oops, something went wrong :(</h2>
        <p className="mb-4">Error message: {error.message}</p>
        <Button onClick={reset} variant="outline">
          Reload
        </Button>
      </main>
    </div>
  );
}
