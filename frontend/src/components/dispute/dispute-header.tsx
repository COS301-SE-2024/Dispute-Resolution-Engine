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
}: {
  id: string;
  label: string;
  startDate: string;
  status: string;
  state: State;
  children?: ReactNode | ReactNode[];
}) {
  return (
    <TooltipProvider>
      <header className="p-4 py-6 grid grid-cols-[1fr_auto]">
        <div>
          <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
            {label}
          </h1>
          <p className="mb-4">Dispute ID: {id}</p>
        </div>

        <dl className="grid grid-cols-2 gap-2">
          <dt className="text-right font-bold">Started:</dt>
          <dd>{startDate}</dd>
          <dt className="text-right font-bold">Status:</dt>
          <dd>{status}</dd>
          <dt className="text-right font-bold">State:</dt>
          <dd className="flex items-center gap-2">
            {state.label}
            <Tooltip>
              <TooltipTrigger>
                <InfoIcon size="1rem" />
              </TooltipTrigger>
              <TooltipContent>{state.description}</TooltipContent>
            </Tooltip>
          </dd>
        </dl>
        <div className="col-span-2 flex gap-2 items-center">{children}</div>
      </header>
    </TooltipProvider>
  );
}
