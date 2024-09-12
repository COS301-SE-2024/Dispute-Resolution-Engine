import { StatusBadge } from "@/components/admin/status-dropdown";
import { TableCell, TableRow } from "@/components/ui/table";
import { AdminDispute } from "@/lib/types/dispute";
import { LinkIcon } from "lucide-react";
import Link from "next/link";

export default function DisputeRow(props: AdminDispute) {
  return (
    <TableRow className="text-nowrap truncate">
      <TableCell className="font-medium">
        <Link href={{ pathname: "/disputes", query: { id: props.id } }}>{props.title}</Link>
      </TableCell>
      <TableCell>
        <StatusBadge variant="active">Hello</StatusBadge>
      </TableCell>
      <TableCell>
        <Link
          className="flex gap-1 items-center hover:underline text-nowrap"
          href={`/workflows/${props.workflow.id}`}
        >
          <span>{props.workflow.title}</span>
          <LinkIcon size="0.8rem" />
        </Link>
      </TableCell>
      <TableCell className="text-center">{props.date_filed}</TableCell>
      <TableCell className="text-center">{props.date_resolved ?? "-"}</TableCell>
    </TableRow>
  );
}
