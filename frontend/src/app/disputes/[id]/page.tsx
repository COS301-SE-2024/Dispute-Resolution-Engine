import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

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
    <header className="p-2 flex">
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

export default function DisputePage({ params }: { params: { id: string } }) {
  return (
    <div className="grow overflow-y-auto flex flex-col">
      <DisputeHeader id={params.id} label="Dispute label" startDate={new Date()} status="Active" />
      <Separator />
      <main className="grow overflow-y-auto p-2"></main>
      <Separator />
      <footer className="p-2 flex justify-between">
        <Button>Action</Button>
        <Button variant="destructive">Action</Button>
      </footer>
    </div>
  );
}
