"use client";

import PageHeader from "@/components/admin/page-header";
import StatusPieChart from "@/components/analytics/dispute-status-pie";
import TicketStatusPieChart from "@/components/analytics/ticket-status-pie";
import { useQuery } from "@tanstack/react-query";
import { QueryProvider } from "./page-client";
import {
  getDisputeCountByStatus,
  getExpertsObjectionSummary,
  getMonthlyDisputes,
  getTicketCountByStatus,
} from "@/lib/api/analytics";
import { useErrorToast } from "@/lib/hooks/use-query-toast";
import { ObjectionBarChart } from "@/components/analytics/objections-bars";
import { MonthlyChart } from "@/components/analytics/monthly-chart";

export default function Home() {
  return (
    <QueryProvider>
      <HomeInner />
    </QueryProvider>
  );
}
function HomeInner() {
  const disputeStatus = useQuery({
    queryKey: ["disputeStatuses"],
    queryFn: () => getDisputeCountByStatus(),
  });
  useErrorToast(disputeStatus.error, "Failed to fetch dispute statistics");

  const ticketStatus = useQuery({
    queryKey: ["ticketStatuses"],
    queryFn: () => getTicketCountByStatus(),
  });
  useErrorToast(ticketStatus.error, "Failed to fetch ticket statistics");

  const expertObjections = useQuery({
    queryKey: ["expertObjections"],
    queryFn: () => getExpertsObjectionSummary(),
  });
  useErrorToast(expertObjections.error, "Failed to fetch objection statistics");

  const monthlyDisputes = useQuery({
    queryKey: ["monthlyDisputes"],
    queryFn: () => getMonthlyDisputes(),
  });

  return (
    <div className="flex flex-col">
      <PageHeader label="Dashboard" />
      <div className="grow md:p-10 md:gap-10 overflow-y-auto grid md:grid-cols-2 grid-cols-1  items-start justify-start">
        {disputeStatus.data && (
          <StatusPieChart
            title="Disputes"
            description="An overview of the disputes created within the last month"
            data={disputeStatus.data}
          />
        )}
        {ticketStatus.data && (
          <TicketStatusPieChart
            title="Tickets"
            description="An overview of the tickets created within the last month"
            data={ticketStatus.data}
          />
        )}
        {expertObjections.data && (
          <ObjectionBarChart
            title="Objections"
            description="How many objections were submitted for each expert"
            data={expertObjections.data}
          />
        )}
        {monthlyDisputes.data && <MonthlyChart data={monthlyDisputes.data} />}
      </div>
    </div>
  );
}
