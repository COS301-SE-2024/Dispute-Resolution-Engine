"use client";

import { Bar, BarChart, CartesianGrid, XAxis } from "recharts";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { useMemo } from "react";

const chartConfig = {} satisfies ChartConfig;

export function ObjectionBarChart({
  title,
  description,
  data,
}: {
  title: string;
  description: string;
  data: Record<string, number>;
}) {
  const processed = useMemo(
    () =>
      Object.entries(data).map(([key, value]) => ({
        name: key,
        count: value,
      })),
    [data]
  );

  return (
    <Card className="mx-0">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart accessibilityLayer data={processed}>
            <CartesianGrid vertical={false} />
            <XAxis dataKey="name" tickLine={false} tickMargin={10} axisLine={false} />
            <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
            <Bar dataKey="count" radius={8} fill="#78A8FF" />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
