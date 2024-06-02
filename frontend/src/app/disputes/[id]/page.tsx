import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
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
  startDate: Date;
  status: string;
}) {
  return (
    <header className="p-4 py-6 flex">
      <div className="grow">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">{label}</h1>
        <span>Started: {startDate.toLocaleDateString()}</span>
      </div>

      <dl className="grid grid-cols-2 gap-2">
        <dt className="text-right font-bold">Dispute ID: </dt>
        <dd>{id}</dd>
        <dt className="text-right font-bold">Status: </dt>
        <dd>{status}</dd>
      </dl>
    </header>
  );
}

export default function DisputePage({ params }: Props) {
  return (
    <div className="grow overflow-y-auto flex flex-col">
      <DisputeHeader id={params.id} label="Dispute label" startDate={new Date()} status="Active" />
      <Separator />
      <ScrollArea className="grow overflow-y-auto p-4">
        <Card className="mb-4">
          <CardHeader>
            <CardTitle>Dispute Details</CardTitle>
          </CardHeader>
          <CardContent>
            <h4>Description</h4>
            <p className="text-sm text-white/70 mt-4">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla pretium, erat non
              iaculis placerat, sem elit aliquet felis, aliquet vehicula ex nunc quis enim. Morbi
              vitae eros ac leo rutrum consectetur. Donec condimentum lobortis ultricies. Donec quam
              mauris, viverra eget lacinia quis, vehicula eu orci. Nunc iaculis iaculis rhoncus.
              Pellentesque blandit enim non elit tincidunt fringilla. Sed eu iaculis libero. Ut nunc
              libero, luctus non dapibus ac, vehicula sed sapien. Donec leo diam, posuere ac dui
              sed, consequat condimentum velit. Aliquam felis sem, iaculis nec consequat eget,
              tincidunt vel justo. Nam id turpis sed neque elementum dictum. Nulla eleifend nibh
              dolor, ac tempus ligula eleifend vitae. Vestibulum dapibus ac libero sit amet aliquam.
              Pellentesque laoreet dolor a orci tristique egestas.
            </p>
            <h4 className="mt-4">Evidence</h4>
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
