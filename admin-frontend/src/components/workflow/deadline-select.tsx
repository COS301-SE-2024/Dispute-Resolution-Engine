"use client";

import { changeDisputeDeadline, getDisputeWorkflow } from "@/lib/api/dispute";
import {} from "@/lib/api/workflow";
import { WORKFLOW_STATES_KEY } from "@/lib/constants";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { useToast } from "@/lib/hooks/use-toast";
import { cn } from "@/lib/utils";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CalendarIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { format } from "date-fns";
import { DateBefore } from "react-day-picker";

export default function DeadlineSelect({ dispute }: { dispute: number }) {
  const query = useQuery({
    queryKey: [WORKFLOW_STATES_KEY, dispute],
    queryFn: () => getDisputeWorkflow(dispute),
  });
  useErrorToast(query.error, "Failed to fetch dispute workflow");

  const { toast } = useToast();
  const client = useQueryClient();
  const currentDeadline = useMutation({
    mutationFn: (state: Date) => changeDisputeDeadline(dispute, state),
    onSuccess: (data, variables) => {
      client.invalidateQueries({
        queryKey: [WORKFLOW_STATES_KEY, dispute],
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

  function onValueChange(value: Date | undefined) {
    if (!value) {
      return;
    }
    if (value.toISOString() === query.data!.current_deadline) {
      return;
    }
    currentDeadline.mutate(value);
  }

  const [date, setDate] = useState<Date>();
  useEffect(() => {
    const { data } = query;
    if (!data) {
      setDate(undefined);
      return;
    }

    if (!data.definition.states[data.current_state].timer) {
      setDate(undefined);
      return;
    }

    if (data.current_deadline) {
      setDate(new Date(data.current_deadline));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query.data]);

  const beforeMatcher: DateBefore = { before: new Date() };
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          disabled={!date}
          variant={"outline"}
          className={cn(
            "w-[280px] bg-surface-light-100 dark:bg-surface-dark-900 justify-start text-left font-normal",
            !date && "text-muted-foreground"
          )}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {date ? (
            format(date, "PPP")
          ) : (
            <span>No deadline {query.isSuccess && !date && "(not supported)"}</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          disabled={beforeMatcher}
          mode="single"
          selected={date}
          onSelect={onValueChange}
          initialFocus
        />
      </PopoverContent>
    </Popover>
  );
}
