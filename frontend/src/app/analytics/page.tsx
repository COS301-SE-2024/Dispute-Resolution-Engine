import { fetchUserCountryDistribution } from "@/lib/api/analytics";
import CountryChart from "./country-chart";

export default async function AnalyticsPage() {
  const chartData = await fetchUserCountryDistribution();
  return (
    <>
      <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">Analytics</h1>
      {chartData.data ? <CountryChart data={chartData.data!} /> : <p>{chartData.error}</p>}
    </>
  );
}
