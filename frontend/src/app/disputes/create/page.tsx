import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Create Dispute",
};

export default function CreateDispute() {
  return (
    <div className="grow overflow-y-auto flex flex-col">
      <header className="px-3 py-6 flex">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
          Create a Dispute
        </h1>
      </header>
      <Separator />
      <main className="grow overflow-y-auto p-2"></main>
      <Separator />
      <footer className="p-2 flex justify-between">
        <Button>Create</Button>
        <Button variant="destructive">Cancel</Button>
      </footer>
    </div>
  );
}
