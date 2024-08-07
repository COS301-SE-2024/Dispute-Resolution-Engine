import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { getDisputeDetails } from "@/lib/api/dispute";
import { Metadata } from "next";

import DisputeClientPage from "./client-page";
import StatusDropdown from "@/app/disputes/[id]/dropdown";
import ExpertRejectForm from "@/components/dispute/expert-reject-form";
import { getAuthToken } from "@/lib/util/jwt";

import { jwtDecode } from "jwt-decode";

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

function DisputeHeader({
  id,
  label,
  startDate,
  status: initialStatus,
}: {
  id: string;
  label: string;
  startDate: string;
  status: string;
}) {
  // TODO: Add contracts for this
  const role = (jwtDecode(getAuthToken()) as any).user.role;

  return (
    <header className="p-4 py-6 flex items-start">
      <div className="grow">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">{label}</h1>
        <p className="mb-4">Started: {startDate}</p>
        {/*TODO: Figure out the conditions for displaying expert rejection */}
        {role == "expert" && <ExpertRejectForm expertId="2" disputeId="x" />}
      </div>

      <dl className="grid grid-cols-2 gap-2">
        <dt className="text-right font-bold">Dispute ID:</dt>
        <dd>{id}</dd>
        <dt className="text-right font-bold">Status:</dt>
        <dd>
          <StatusDropdown disputeId={id} status={initialStatus} />
        </dd>
      </dl>
    </header>
  );
}
