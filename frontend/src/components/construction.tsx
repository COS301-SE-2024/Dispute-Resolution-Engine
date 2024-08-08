import { BackButton } from "@/components/back-button";
import { Construction } from "lucide-react";

export default function UnderConstruction() {
  return (
    <main className="flex flex-col items-center justify-center h-full">
      <Construction size="3rem" className="mb-2" />
      <h1 className="font-bold text-center text-3xl dark:text-white text-black mb-2">
        Under Construction
      </h1>
      <p className="mb-10 dark:text-dre-bg-light/50">
        This page is under construction. Watch this space!
      </p>
      <BackButton>Back</BackButton>
    </main>
  );
}
