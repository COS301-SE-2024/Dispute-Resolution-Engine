"use client";

import { TrendingUp } from "lucide-react";
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts";

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
import { useMemo } from "react";

export const description = "A simple area chart";

const chartData = [
  { month: "January", desktop: 186 },
  { month: "February", desktop: 305 },
  { month: "March", desktop: 237 },
  { month: "April", desktop: 73 },
  { month: "May", desktop: 209 },
  { month: "June", desktop: 214 },
];

const chartConfig = {
  desktop: {
    label: "Desktop",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig;

export function MonthlyChart({ data }: { data: Record<string, number> }) {
  const processed = useMemo(
    () =>
      Object.keys(data)
        .sort()
        .map((key) => ({
          time: key,
          count: data[key],
        })),
    [data]
  );

  return (
    <Card className="mx-0">
      <CardHeader>
        <CardTitle>Dispute Length</CardTitle>
        <CardDescription>Showing how long disputes take over time</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <AreaChart
            accessibilityLayer
            data={processed}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis dataKey="time" tickLine={false} axisLine={false} tickMargin={8} />
            <ChartTooltip cursor={false} content={<ChartTooltipContent indicator="line" />} />
            <Area
              dataKey="count"
              type="natural"
              fill="#78A8FF"
              fillOpacity={0.4}
              stroke="#78A8FF"
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
