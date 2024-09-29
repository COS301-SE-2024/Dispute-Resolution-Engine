"use client";

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectLabel,
  SelectItem,
} from "@/components/ui/select";
import { DISPUTE_STATUS, DisputeStatus } from "@/lib/types";

export default function StatusFilter({
  initialValue,
  onValueChange = () => {},
}: {
  onValueChange?: (status: DisputeStatus | undefined) => void;
  initialValue?: DisputeStatus;
}) {
  return (
    <Select
      defaultValue={initialValue ?? "none"}
      onValueChange={(val) => onValueChange(val === "none" ? undefined : (val as DisputeStatus))}
    >
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="No status" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Dispute status</SelectLabel>
          <SelectItem value={"none"}>No status</SelectItem>
          {DISPUTE_STATUS.map((status) => (
            <SelectItem key={status} value={status}>
              {status}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
