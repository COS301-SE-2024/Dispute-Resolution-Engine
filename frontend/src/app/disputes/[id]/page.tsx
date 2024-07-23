import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { getDisputeDetails } from "@/lib/api/dispute";
import { Metadata } from "next";

import DisputeHeader from "@/app/disputes/[id]/Dropdown";
import DisputeClientPage from "./client-page";

type Props = {
  params: { id: string };
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  return {
    title: `Dispute ${params.id}`,
    description: "Dispute description",
  };
}

export default async function DisputePage({ params }: Props) {
  const { data, error } = await getDisputeDetails(params.id);
  if (error || !data) {
    return <h1>{error}</h1>;
  }

  return (
    <div className="grow overflow-y-auto flex flex-col">
      <DisputeHeader
        id={data.id}
        label={data.title}
        startDate={data.case_date.substring(0, 10)}
        status={data.status}
      />
      <Separator />
      <ScrollArea className="grow overflow-y-auto p-4">
        <DisputeClientPage data={data} />
      </ScrollArea>
      <Separator />
    </div>
  );
}
