"use client";

import * as React from "react";
import { TrendingUp } from "lucide-react";
import { Label, Legend, Pie, PieChart } from "recharts";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { TICKET_STATUS, TicketStatus } from "@/lib/types/tickets";

export const description = "A donut chart with text";

const ACTIVE_STATUS: TicketStatus[] = ["Open", "On Hold"];

const CHART_COLORS: Record<TicketStatus, string> = {
  Open: "#78A8FF",
  "On Hold": "#fb923c",

  Closed: "#22c55e",
  Solved: "#22c55e",
};

const chartConfig = {} satisfies ChartConfig;

export default function TicketStatusPieChart({
  title,
  description,
  data,
}: {
  title: string;
  description: string;
  data: Record<TicketStatus, number>;
}) {
  const totalCount = React.useMemo(() => {
    return Object.values(data).reduce((acc, curr) => acc + curr, 0);
  }, [data]);

  const totalOpen = React.useMemo(() => {
    const entries = Object.entries(data) as [TicketStatus, number][];
    return entries.reduce((acc, curr) => {
      if (ACTIVE_STATUS.includes(curr[0])) {
        return acc + curr[1];
      }
      return acc;
    }, 0);
  }, [data]);

  const processed = TICKET_STATUS.map((status) => ({
    status,
    count: data[status],
    fill: CHART_COLORS[status],
  }));

  return (
    <Card className="mx-0">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="mx-auto aspect-square max-h-[250px]">
          <PieChart>
            <ChartTooltip cursor={false} content={<ChartTooltipContent />} />
            <Pie
              data={processed}
              dataKey="count"
              nameKey="status"
              innerRadius={60}
              paddingAngle={2}
            >
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={viewBox.cy}
                          className="dark:fill-white text-3xl font-bold"
                        >
                          {totalCount.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="dark:fill-white/50"
                        >
                          Tickets
                        </tspan>
                      </text>
                    );
                  }
                }}
              />
            </Pie>
          </PieChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex flex-col items-stretch">
        <p className="flex">
          <strong className="grow">Open tickets:</strong>
          <span>{((totalOpen / totalCount) * 100).toFixed(0)}%</span>
        </p>
        <p className="flex">
          <strong className="grow">Closed tickets:</strong>
          <span>{((1 - totalOpen / totalCount) * 100).toFixed(0)}%</span>
        </p>
      </CardFooter>
    </Card>
  );
}
