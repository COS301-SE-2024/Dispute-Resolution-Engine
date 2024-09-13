import Link from "next/link";
import { Button } from "@/components/ui/button";
import Navbar from "@/components/navbar";

export default function NotFound() {
  return (
    <div className="grid grid-cols-1 grid-rows-[auto_1fr] h-full overflow-hidden">
      <Navbar />
      <div className="flex flex-col items-center justify-center h-full">
        <main>
          <h2 className="font-bold text-3xl mb-2">404 Page not found</h2>
          <p className="mb-4">We could not find what you are looking for...</p>
          <Button asChild variant="outline">
            <Link href="/">Home</Link>
          </Button>
        </main>
      </div>
    </div>
  );
}
