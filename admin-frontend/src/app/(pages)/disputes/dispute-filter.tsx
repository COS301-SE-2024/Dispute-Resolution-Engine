"use client";

import StatusFilter from "@/components/dispute/status-filter";
import WorkflowFilter from "@/components/dispute/workflow-filter";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { DisputeStatus, Filter } from "@/lib/types";
import { ReactNode, useState } from "react";

export default function DisputeFilter({
  children,
  onValueChange = () => {},
}: {
  children: ReactNode;
  onValueChange?: (filter: Filter[]) => void;
}) {
  const [status, setStatus] = useState<DisputeStatus | undefined>(undefined);
  const [workflow, setWorkflow] = useState<string | undefined>(undefined);

  function applyFilter() {
    let filter: Filter[] = [];
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
        <StatusFilter onValueChange={setStatus} />
        <label>Workflow</label>
        <WorkflowFilter onValueChange={setWorkflow} />
        <div className="col-span-2 flex flex-end">
          <Button className="ml-auto" onClick={applyFilter}>
            Apply
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}
