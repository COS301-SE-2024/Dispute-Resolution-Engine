import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import dynamic from "next/dynamic";
const CreateDisputeClient = dynamic(() => import("@/app/disputes/create/CreateDisputeClient"), {
  ssr: false,
});

export default function CreateDispute() {
  return (
    <div className="grow overflow-y-auto flex flex-col">
      <header className="p-4 py-6 flex">
        <div className="grow">
          <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
            Create a dispute
          </h1>
        </div>
      </header>
      <Separator />
      <ScrollArea>
        <CreateDisputeClient />
      </ScrollArea>
    </div>
  );
}
