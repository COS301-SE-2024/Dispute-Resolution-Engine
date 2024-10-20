"use client";

import StatusFilter from "@/components/dispute/status-filter";
import WorkflowFilter from "@/components/dispute/workflow-filter";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { DisputeFilter, DisputeStatus, Filter } from "@/lib/types";
import { ReactNode, useState } from "react";

export default function DisputeFilters({
  children,
  value = [],
  onValueChange = () => {},
}: {
  children: ReactNode;
  onValueChange?: (filter: DisputeFilter[]) => void;
  value?: DisputeFilter[];
}) {
  const statusFilter = value.find((f) => f.attr == "status")?.value as DisputeStatus;
  const workflowFilter = value.find((f) => f.attr == "workflow")?.value;

  const [status, setStatus] = useState<DisputeStatus | undefined>(undefined);
  const [workflow, setWorkflow] = useState<string | undefined>(undefined);

  function applyFilter() {
    let filter: DisputeFilter[] = [];
    if (status) {
      filter.push({
        attr: "status",
        value: status,
      });
    }
    if (workflow) {
      filter.push({
        attr: "workflow",
        value: workflow,
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
        <label>Workflow</label>
        <WorkflowFilter initialValue={workflowFilter} onValueChange={setWorkflow} />
        <div className="col-span-2 flex flex-end">
          <Button className="ml-auto" onClick={applyFilter}>
            Apply
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}
