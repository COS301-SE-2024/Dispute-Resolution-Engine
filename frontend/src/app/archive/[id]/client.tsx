"use client";

import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { useRouter } from "next/navigation";

export function BackButton() {
  const router = useRouter();

  return (
    <Button
      variant="ghost"
      onClick={() => router.back()}
      className="rounded-full aspect-square p-0 flex justify-center"
      aria-label="Back"
    >
      <ArrowLeft />
    </Button>
  );
}
