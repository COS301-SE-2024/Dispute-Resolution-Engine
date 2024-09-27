"use client";

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectItem,
} from "@/components/ui/select";
import { getWorkflowList } from "@/lib/api/workflow";
import { SelectLabel, SelectProps } from "@radix-ui/react-select";
import { useQuery } from "@tanstack/react-query";

export default function WorkflowSelect({
  id,
  ...props
}: SelectProps & {
  id?: string;
}) {
  const query = useQuery({
    queryKey: ["workflowList"],
    queryFn: () => getWorkflowList({}),
    staleTime: Infinity,
  });

  return (
    <Select {...props}>
      <SelectTrigger disabled={query.isPending}>
        <SelectValue id={id} placeholder="Select a workflow" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Workflows</SelectLabel>
          {query.data?.workflows.map((wf) => (
            <SelectItem key={wf.id} value={wf.name}>
              {wf.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
