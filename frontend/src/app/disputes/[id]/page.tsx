import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { getDisputeDetails, updateDisputeStatus } from "@/lib/api/dispute";
import { Metadata } from "next";
import { File, WorkflowIcon } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
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
  const [funcStatus, setStatus] = useState("")
  const handleStatusChange = async (newStatus: string) => {
    const result = await updateDisputeStatus(id, newStatus);
    if (!result.error) {
      setStatus(newStatus); // Update status state if API call is successful
    } else {
      console.error("Failed to update status:", result.error);
    }
  }
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
          <DropdownMenu >
            <DropdownMenuTrigger asChild>
              <Badge>{status}</Badge>
              </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>Next Steps</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Waiting for Admin Approval</DropdownMenuItem>
              <DropdownMenuItem>Waiting for respondent</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
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
  var eviStr : string = "No evidence provided"
  if(data.evidence){
    if(data.evidence[0]){
      var eviStr : string = data.evidence[0].toString()
      eviStr = eviStr.split("/").pop() as string
    }
  }


  return (
    <div className="grow overflow-y-auto flex flex-col">
      <DisputeHeader
        id={data.id}
        label={data.title}
        startDate={data.case_date.substring(0, 10)}
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
            <CardTitle>Complainant&apos;s Evidence</CardTitle>

          </CardHeader>
          <CardContent>
            <div className="rounded-lg bg-gray-950 p-4 text-center text-gray-50 w-40">
              <File className="mx-auto h-8 w-8" />
              <p className="mt-2 text-sm font-medium">{eviStr}</p>
            </div>
          </CardContent>
        </Card>
      </ScrollArea>
      <Separator />
    </div>
  );
}
