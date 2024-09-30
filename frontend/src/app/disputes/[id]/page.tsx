import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { getDisputeDetails, getDisputeWorkflow } from "@/lib/api/dispute";
import { Metadata } from "next";

import DisputeClientPage from "./client-page";
import StatusDropdown from "@/app/disputes/[id]/dropdown";
import ExpertRejectForm from "@/components/dispute/expert-reject-form";
import { getAuthToken } from "@/lib/util/jwt";

import { jwtDecode } from "jwt-decode";
import DisputeDecisionForm from "@/components/dispute/decision-form";
import CreateTicketDialog from "@/components/dispute/ticket-form";
import { Button } from "@/components/ui/button";
import WorkflowSelect from "@/components/form/workflow-select";
import DisputeHeader from "@/components/dispute/dispute-header";
import Link from "next/link";
import { State } from "@/lib/interfaces/workflow";

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
  const workflow = await getDisputeWorkflow(params.id);
  if (error || !data) {
    return <h1>{error}</h1>;
  }

  return (
    <div className="grow overflow-y-auto flex flex-col">
      <DisputeHeader2
        id={data.id}
        label={data.title}
        startDate={data.case_date.substring(0, 10)}
        status={data.status}
        state={workflow.definition.states[workflow.current_state]}
      />
      <Separator />
      <ScrollArea className="grow overflow-y-auto p-4">
        <DisputeClientPage data={data} />
      </ScrollArea>
      <Separator />
    </div>
  );
}

function DisputeHeader2(props: {
  id: string;
  label: string;
  startDate: string;
  status: string;
  state: State;
}) {
  // TODO: Add contracts for this
  const user = (jwtDecode(getAuthToken()) as any).user.id;
  const role = (jwtDecode(getAuthToken()) as any).user.role;

  return (
    <DisputeHeader {...props}>
      {role == "expert" && <ExpertRejectForm expertId={user} disputeId={props.id} />}

      {role == "expert" && (
        <DisputeDecisionForm disputeId={props.id} asChild>
          <Button>Render decision</Button>
        </DisputeDecisionForm>
      )}

      <Button variant="outline" asChild>
        <Link href={`/disputes/${props.id}/tickets`}>Go to tickets</Link>
      </Button>
    </DisputeHeader>
  );
}
