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

export default function WorkflowFilter({
  onValueChange = () => {},
}: {
  onValueChange?: (id: string | undefined) => void;
}) {
  return (
    <Select
      defaultValue="none"
      onValueChange={(val) => onValueChange(val === "none" ? undefined : val)}
    >
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="No workflow" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Workflow</SelectLabel>
          <SelectItem value={"none"}>No workflow</SelectItem>
          <SelectItem value={"1"}>Workflow #1</SelectItem>
          <SelectItem value={"2"}>Workflow #2</SelectItem>
          <SelectItem value={"3"}>Workflow #3</SelectItem>
          <SelectItem value={"4"}>Workflow #4</SelectItem>
          <SelectItem value={"5"}>Workflow #5</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
