import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { getDisputeDetails } from "@/lib/api/dispute";
import { Metadata } from "next";

type Props = {
  params: { id: string };
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  return {
    title: `Dispute ${params.id}`,
    description: "Dispute description",
  };
}

function DisputeHeader({
  id,
  label,
  startDate,
  status,
}: {
  id: string;
  label: string;
  startDate: string;
  status: string;
}) {
  return (
    <header className="p-4 py-6 flex">
      <div className="grow">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">{label}</h1>
        <span>Started: {startDate}</span>
      </div>

      <dl className="grid grid-cols-2 gap-2">
        <dt className="text-right font-bold">Dispute ID: </dt>
        <dd>{id}</dd>
        <dt className="text-right font-bold">Status: </dt>
        <dd>
          <Badge>{status}</Badge>
        </dd>
      </dl>
    </header>
  );
}

export default async function DisputePage({ params }: Props) {
  const { data, error } = await getDisputeDetails(params.id);

  if (error || !data) {
    return <h1>{error}</h1>;
  }

  return (
    <div className="grow overflow-y-auto flex flex-col">
      <DisputeHeader
        id={data.id}
        label={data.title}
        startDate={data.date_created}
        status={data.status}
      />
      <Separator />
      <ScrollArea className="grow overflow-y-auto p-4">
        <Card className="mb-4">
          <CardHeader>
            <CardTitle>Description</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-white/70 mt-4">{data.description}</p>
          </CardContent>
        </Card>
        <Card className="mb-4">
          <CardHeader>
            <CardTitle>Respondant&apos;s Evidence</CardTitle>
            <CardDescription>All the evidence the respondant has submitted</CardDescription>
          </CardHeader>
        </Card>
        <Card className="mb-4">
          <CardHeader>
            <CardTitle>Respondant Information</CardTitle>
            <CardDescription>Who you gon&apos; sue?</CardDescription>
          </CardHeader>
        </Card>
      </ScrollArea>
      <Separator />
    </div>
  );
}
