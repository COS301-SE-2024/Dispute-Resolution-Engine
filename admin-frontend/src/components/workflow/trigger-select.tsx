"use client";

import { getWorkflowTriggers } from "@/lib/api/workflow";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { SelectProps } from "@radix-ui/react-select";
import {
  SelectGroup,
  SelectItem,
  SelectLabel,
  Select,
  SelectContent,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { useQuery } from "@tanstack/react-query";

type TriggerSelectProps = Omit<SelectProps, "disabled"> & {
  id: string;
};

export default function TriggerSelect({ id, ...props }: TriggerSelectProps) {
  const query = useQuery({
    queryKey: ["triggers"],
    queryFn: () => getWorkflowTriggers(),
  });
  useErrorToast(query.error, "Failed to fetch dispute workflow");

  return (
    <Select disabled={!query.isSuccess} {...props}>
      <SelectTrigger id={id} className="w-[180px]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Triggers</SelectLabel>
          {query.data?.map((id) => (
            <SelectItem key={id} value={id}>
              {id}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
