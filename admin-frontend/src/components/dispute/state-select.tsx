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
import { changeDisputeState, getDisputeWorkflow } from "@/lib/api/dispute";
import {} from "@/lib/api/workflow";
import { WORKFLOW_STATES_KEY } from "@/lib/constants";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { useToast } from "@/lib/hooks/use-toast";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useMemo } from "react";

export default function StateSelect({ dispute }: { dispute: number }) {
  const query = useQuery({
    queryKey: [WORKFLOW_STATES_KEY, dispute],
    queryFn: () => getDisputeWorkflow(dispute),
  });
  useErrorToast(query.error, "Failed to fetch dispute workflow");

  const states = useMemo(() => {
    if (query.data) {
      return Object.entries(query.data.definition.states).map(([id, st]) => ({
        id,
        label: st.label,
      }));
    }
    return undefined;
  }, [query.data]);

  const { toast } = useToast();
  const client = useQueryClient();
  const currentState = useMutation({
    mutationFn: (state: string) => changeDisputeState(dispute, state),
    onSuccess: (data, variables) => {
      client.invalidateQueries({
        queryKey: [WORKFLOW_STATES_KEY, dispute],
      });
      toast({
        title: "Status updated successfully",
      });
    },
    onError: (error) => {
      toast({
        variant: "error",
        title: "Something went wrong",
        description: error?.message,
      });
    },
  });

  function onValueChange(value: string) {
    if (value === query.data!.current_state) {
      return;
    }
    currentState.mutate(value);
  }

  return (
    <Select
      disabled={!query.isSuccess}
      defaultValue={query.data?.current_state}
      onValueChange={onValueChange}
    >
      <SelectTrigger className="w-[180px]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Workflow state</SelectLabel>
          {states?.map(({ id, label }) => (
            <SelectItem key={id} value={id}>
              {label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
