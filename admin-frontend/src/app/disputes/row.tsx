"use client";

import { StatusBadge, StatusDropdown } from "@/components/admin/status-dropdown";
import { TableCell, TableRow } from "@/components/ui/table";
import { AdminDispute } from "@/lib/types/dispute";
import Link from "next/link";

export default function DisputeRow(props: AdminDispute) {
  function changeStatus(status: string) {
    if (props.status === status) {
      return;
    }

    // TODO: Integrate with API
    alert(`Change dispute status to ${status}`);
  }

  return (
    <TableRow>
      <TableCell className="font-medium">
        <Link href={{ pathname: "/disputes", query: { id: props.id } }}>{props.title}</Link>
      </TableCell>
      <TableCell>
        <StatusDropdown onSelect={changeStatus}>
          <StatusBadge variant="active" dropdown>
            Hello
          </StatusBadge>
        </StatusDropdown>
      </TableCell>
      <TableCell>
        <Link href={`/workflows/${props.workflow.id}`}>{props.workflow.title}</Link>
      </TableCell>
      <TableCell className="text-center">{props.date_filed}</TableCell>
      <TableCell className="text-center">{props.date_resolved ?? "-"}</TableCell>
    </TableRow>
  );
}
