import { getDisputeDetails, getDisputeWorkflow } from "@/lib/api/dispute";
import { Metadata } from "next";
import DisputeClientPage from "./client-page";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { State } from "@/lib/interfaces/workflow";
import {
  BadgeCheckIcon,
  ChevronLeftIcon,
  EllipsisVerticalIcon,
  InfoIcon,
  TicketCheckIcon,
  TriangleAlertIcon,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getAuthToken } from "@/lib/util/jwt";
import ExpertRejectForm from "@/components/dispute/expert-reject-form";
import DisputeDecisionForm from "@/components/dispute/decision-form";

import { jwtDecode } from "jwt-decode";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";

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
    <div className="grid grid-rows-[auto_1fr] overflow-y-hidden">
      <DisputeHeader
        id={data.id}
        title={data.title}
        date={data.case_date}
        status={data.status}
        state={workflow.definition.states[workflow.current_state]}
      />
      <div className="overflow-y-auto p-5">
        <DisputeClientPage data={data} />
      </div>
    </div>
  );
}

function DisputeHeader(props: {
  id: string;
  title: string;
  date: string;
  status: string;
  state: State;
}) {
  const user = (jwtDecode(getAuthToken()) as any).user.id;
  const role = (jwtDecode(getAuthToken()) as any).user.role;

  return (
    <div className="border-b border-primary-500/30 grid grid-cols-[auto_1fr_auto] md:grid-cols-[auto_1fr_auto_auto] p-3 gap-3">
      <Button asChild className="p-2 rounded-full" variant="ghost">
        <Link href="/disputes">
          <ChevronLeftIcon />
        </Link>
      </Button>

      <section className="mt-1">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl text-wrap">
          {props.title}
        </h1>
        <p>Filed: {props.date}</p>
      </section>

      <dl className="col-start-2 md:col-start-3 mt-1 grid grid-cols-[auto_1fr] gap-x-2">
        <dt className="font-bold">Status:</dt>
        <dd>{props.status}</dd>
        <dt className="font-bold">State:</dt>
        {/* <Tooltip>
            <TooltipTrigger asChild> */}
        <dd className="flex items-center gap-2">
          {props.state.label}
          <InfoIcon size="1rem" />
        </dd>
        {/* </TooltipTrigger>
            <TooltipContent>{state.description}</TooltipContent>
          </Tooltip> */}
      </dl>

      <Popover>
        <PopoverTrigger asChild>
          <Button className="p-2 rounded-full row-start-1 col-start-4" variant="ghost">
            <EllipsisVerticalIcon />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="flex flex-col p-2">
          <Button className="gap-2 p-2" variant="ghost" asChild>
            <Link href={`/disputes/${props.id}/tickets`}>
              <TicketCheckIcon />
              <span>Dispute tickets</span>
            </Link>
          </Button>
          {role == "expert" && (
            <>
              <DisputeDecisionForm disputeId={props.id} asChild>
                <Button className="gap-2 p-2" variant="ghost">
                  <BadgeCheckIcon />
                  <span>Render decision</span>
                </Button>
              </DisputeDecisionForm>
              <ExpertRejectForm expertId={user} disputeId={props.id} asChild>
                <Button className="gap-2 p-2 text-red-500" variant="ghost">
                  <TriangleAlertIcon />
                  <span>Object to assignment</span>
                </Button>
              </ExpertRejectForm>
            </>
          )}
        </PopoverContent>
      </Popover>
    </div>
  );
}

// function DisputeHeader2(props: {
//   id: string;
//   label: string;
//   startDate: string;
//   status: string;
//   state: State;
// }) {
//   // TODO: Add contracts for this

//   return (
//     <DisputeHeader {...props}>
//       {role == "expert" && <ExpertRejectForm expertId={user} disputeId={props.id} />}

//       {role == "expert" && (
//         <DisputeDecisionForm disputeId={props.id} asChild>
//           <Button>Render decision</Button>
//         </DisputeDecisionForm>
//       )}

//       <Button variant="outline" asChild>
//         <Link href={`/disputes/${props.id}/tickets`}>Go to tickets</Link>
//       </Button>
//     </DisputeHeader>
//   );
// }
