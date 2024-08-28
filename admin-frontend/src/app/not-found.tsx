import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center">
      <main>
        <h2 className="font-bold text-3xl mb-2">404 Page not found</h2>
        <p className="mb-4">We could not find what you are looking for...</p>
        <Button asChild variant="outline">
          <Link href="/">Back to Dashboard</Link>
        </Button>
      </main>
    </div>
  );
}
