"use client";

import { type ChartConfig, ChartContainer } from "@/components/ui/chart";
import { CountryCount } from "@/lib/api/analytics";
import { Bar, BarChart, XAxis } from "recharts";

const chartConfig = {
  count: {
    color: "#2563eb",
  },
} satisfies ChartConfig;

export default function CountryChart({ data }: { data: CountryCount[] }) {
  return (
    <ChartContainer config={chartConfig} className="min-h-[200px] w-1/2">
      <BarChart accessibilityLayer data={data}>
        <XAxis dataKey="nationality" tickLine={false} tickMargin={10} axisLine={false} />
        <Bar dataKey="count" fill="var(--color-count)" radius={4} />
      </BarChart>
    </ChartContainer>
  );
}
