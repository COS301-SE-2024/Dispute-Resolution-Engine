import PageHeader from "@/components/admin/page-header";
import StatusPieChart from "@/components/analytics/status-pie";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { DisputeStatus } from "@/lib/types";

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

export default function Home() {
  return (
    <div className="flex flex-col">
      <PageHeader label="Dashboard" />
      <div className="grow md:p-10 md:gap-10 overflow-y-auto flex flex-wrap  items-start justify-start">
        <StatusPieChart data={data} />
        <StatusPieChart data={data} />
        <StatusPieChart data={data} />
        <StatusPieChart data={data} />
      </div>
    </div>
  );
}
