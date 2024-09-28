import PageHeader from "@/components/admin/page-header";
import StatusPieChart from "@/components/analytics/dispute-status-pie";
import TicketStatusPieChart from "@/components/analytics/ticket-status-pie";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { DisputeStatus } from "@/lib/types";
import { TicketStatus } from "@/lib/types/tickets";

const data: Record<DisputeStatus, number> = {
  "Awaiting Respondant": 1,
  Active: 2,
  Review: 3,
  Settled: 4,
  Refused: 5,
  Withdrawn: 6,
  Transfer: 7,
  Appeal: 8,
  Other: 9,
};

const data2: Record<TicketStatus, number> = {
  Open: 1,
  Closed: 2,
  Solved: 3,
  "On Hold": 4,
};

export default function Home() {
  return (
    <div className="flex flex-col">
      <PageHeader label="Dashboard" />
      <div className="grow md:p-10 md:gap-10 overflow-y-auto flex flex-wrap  items-start justify-start">
        <StatusPieChart data={data} />
        <TicketStatusPieChart data={data2} />
      </div>
    </div>
  );
}
