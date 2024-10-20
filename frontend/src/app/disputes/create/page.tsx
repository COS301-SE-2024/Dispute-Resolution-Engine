import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import dynamic from "next/dynamic";
import { Content, Header, Root } from "../custom-layout";
const CreateDisputeClient = dynamic(() => import("@/app/disputes/create/CreateDisputeClient"), {
  ssr: false,
});

export default function CreateDispute() {
  return (
    <Root>
      <Header>
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
          Create a dispute
        </h1>
      </Header>
      <Content>
        <CreateDisputeClient />
      </Content>
    </Root>
  );
}
