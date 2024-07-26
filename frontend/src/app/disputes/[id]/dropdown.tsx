"use client"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import { getStatusEnum, updateDisputeStatus } from "@/lib/api/dispute";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";

export default async function DisputeHeader({ id, label, startDate, status: initialStatus }: {
  id: string;
  label: string;
  startDate: string;
  status: string;
}) {
  const [status, setStatus] = useState(initialStatus);
  const handleStatusChange = async (newStatus: string) => {
    try {
      const response = await updateDisputeStatus(id, newStatus);
      console.log("RESPONSE", response)
      setStatus(newStatus)
    } catch (error) {
      console.error("Failed to update dispute status:", error);
    }
  };
  const nextStates = await getStatusEnum()
  const optionsJSX = nextStates.map((state : string, i : number) => (
    <DropdownMenuItem key={i} onSelect={() => handleStatusChange("Waiting for admin approval")}>
      {state}
    </DropdownMenuItem>
  ))

  return (
    <header className="p-4 py-6 flex">
      <div className="grow">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">{label}</h1>
        <span>Started: {startDate}</span>
      </div>

      <dl className="grid grid-cols-2 gap-2">
        <dt className="text-right font-bold">Dispute ID:</dt>
        <dd>{id}</dd>
        <dt className="text-right font-bold">Status:</dt>
        <dd>
          <DropdownMenu>
            <DropdownMenuTrigger>
              {status}
              {/*<Badge>{status}</Badge>*/}
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>Next Steps</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {optionsJSX}
            </DropdownMenuContent>
          </DropdownMenu>
        </dd>
      </dl>
    </header>
  );
}