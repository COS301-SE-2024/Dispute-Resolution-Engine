"use client";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { getStatusEnum, updateDisputeStatus } from "@/lib/api/dispute";
import { useEffect, useState } from "react";

export default function StatusDropdown({
  disputeId,
  status,
}: {
  disputeId: string;
  status: string;
}) {
  const [currentStatus, setCurrentStatus] = useState(status);
  const [nextStates, setNextStates] = useState<string[]>([]);

  useEffect(() => {
    const fetchStatusEnum = async () => {
      try {
        const states = await getStatusEnum();
        setNextStates(states);
      } catch (error) {
        console.error("Failed to fetch status enum:", error);
      }
    };
    fetchStatusEnum();
  }, []);

  const handleStatusChange = async (newStatus: string) => {
    try {
      const response = await updateDisputeStatus(disputeId, newStatus);
      console.log("RESPONSE", response);
      setCurrentStatus(newStatus);
    } catch (error) {
      console.error("Failed to update dispute status:", error);
    }
  };

  const optionsJSX = nextStates.map((state: string, i: number) => (
    <DropdownMenuItem key={i} onSelect={() => handleStatusChange(`${state}`)}>
      {state}
    </DropdownMenuItem>
  ));

  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        {currentStatus}
        {/*<Badge>{status}</Badge>*/}
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>Next Steps</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {optionsJSX}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

