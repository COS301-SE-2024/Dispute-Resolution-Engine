import { State } from "@/lib/interfaces/workflow";
import { type ReactNode } from "react";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip";
import { InfoIcon } from "lucide-react";

export default function DisputeHeader({
  id,
  label,
  startDate,
  status,
  children,
  state,
  estimate,
}: {
  id: string;
  label: string;
  startDate: string;
  status: string;
  state: State;
  children?: ReactNode | ReactNode[];
  estimate: string;
}) {
  return (
    <TooltipProvider>
      <header className="p-4 py-6 grid grid-cols-[1fr_auto] items-start">
        <div>
          <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
            {label}
          </h1>
          <p className="mb-4">Dispute ID: {id}</p>
          <div className="flex gap-2 items-center">{children}</div>
        </div>

        <dl className="grid grid-cols-2 gap-2 items-center">
          <dt className="text-right font-bold">Started:</dt>
          <dd>{startDate}</dd>
          <dt className="text-right font-bold">Status:</dt>
          <dd>{status}</dd>
          <dt className="text-right font-bold">State:</dt>
          <Tooltip>
            <TooltipTrigger asChild>
              <dd className="flex items-center gap-2">
                {state.label}
                <InfoIcon size="1rem" />
              </dd>
            </TooltipTrigger>
            <TooltipContent>{state.description}</TooltipContent>
          </Tooltip>
          <dt className="text-right font-bold">Estimated time:</dt>
          <dd>{estimate}</dd>
        </dl>
      </header>
    </TooltipProvider>
  );
}
