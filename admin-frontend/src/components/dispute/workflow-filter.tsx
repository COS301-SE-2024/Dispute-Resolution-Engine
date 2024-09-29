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
import { getWorkflowList } from "@/lib/api/workflow";
import { WORKFLOW_LIST_KEY } from "@/lib/constants";
import { useQuery } from "@tanstack/react-query";

export default function WorkflowFilter({
  initialValue,
  onValueChange = () => {},
}: {
  onValueChange?: (id: string | undefined) => void;
  initialValue?: string;
}) {
  const query = useQuery({
    queryKey: [WORKFLOW_LIST_KEY],
    queryFn: async () => {
      const { workflows } = await getWorkflowList({});
      return workflows;
    },
  });

  return (
    <Select
      defaultValue={initialValue ?? "none"}
      onValueChange={(val) => onValueChange(val === "none" ? undefined : val)}
    >
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="No workflow" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Workflow</SelectLabel>
          <SelectItem value={"none"}>No workflow</SelectItem>
          {query.data?.map((wf) => (
            <SelectItem key={wf.id} value={wf.id.toString()}>
              {wf.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
