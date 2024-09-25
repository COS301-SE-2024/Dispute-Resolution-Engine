"use client";

import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { TICKET_STATUS, type TicketFilter, TicketStatus } from "@/lib/types/tickets";
import { ReactNode, useState } from "react";

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectLabel,
  SelectItem,
} from "@/components/ui/select";

export default function TicketFilters({
  children,
  value = [],
  onValueChange = () => {},
}: {
  children: ReactNode;
  onValueChange?: (filter: TicketFilter[]) => void;
  value?: TicketFilter[];
}) {
  const statusFilter = value.find((f) => f.attr == "status")?.value as TicketStatus;

  const [status, setStatus] = useState<TicketStatus | undefined>(undefined);

  function applyFilter() {
    let filter: TicketFilter[] = [];
    if (status) {
      filter.push({
        attr: "status",
        value: status,
      });
    }
    onValueChange(filter);
  }

  return (
    <Popover>
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent className="grid gap-x-2 gap-y-3 grid-cols-[auto_1fr] items-center">
        <strong className="col-span-2">Filter</strong>

        <label>Status</label>
        <StatusFilter initialValue={statusFilter} onValueChange={setStatus} />
        <div className="col-span-2 flex flex-end">
          <Button className="ml-auto" onClick={applyFilter}>
            Apply
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}

function StatusFilter({
  initialValue,
  onValueChange = () => {},
}: {
  onValueChange?: (status: TicketStatus | undefined) => void;
  initialValue?: TicketStatus;
}) {
  return (
    <Select
      defaultValue={initialValue ?? "none"}
      onValueChange={(val) => onValueChange(val === "none" ? undefined : (val as TicketStatus))}
    >
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="No filter" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Ticket status</SelectLabel>
          <SelectItem value={"none"}>No status</SelectItem>
          {TICKET_STATUS.map((status) => (
            <SelectItem key={status} value={status}>
              {status}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
