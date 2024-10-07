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
      <header className="p-4 py-6 grid grid-cols-1 md:grid-cols-[1fr_auto] items-start gap-y-4">
        <div>
          <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl text-wrap">
            {label}
          </h1>

          <p>Started: {startDate}</p>
        </div>

        <dl className="grid grid-cols-[auto_1fr] gap-x-2">
          <dt className="font-bold">Status:</dt>
          <dd>{status}</dd>
          <dt className="font-bold">State:</dt>
          <Tooltip>
            <TooltipTrigger asChild>
              <dd className="flex items-center gap-2">
                {state.label}
                <InfoIcon size="1rem" />
              </dd>
            </TooltipTrigger>
            <TooltipContent>{state.description}</TooltipContent>
          </Tooltip>
        </dl>

        <div className="md:col-span-2 flex gap-2 items-center">{children}</div>
      </header>
    </TooltipProvider>
  );
}
