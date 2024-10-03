import DisputeHeader from "@/components/dispute/dispute-header";
import CreateTicketDialog from "@/components/dispute/ticket-form";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getTicketSummaries } from "@/lib/api/tickets";
import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react";
import Link from "next/link";

type Props = {
  params: { id: string };
};

export default async function TicketsPage({ params: { id } }: Props) {
  const data = await getTicketSummaries(parseInt(id));
  return (
    <div className="grid grid-rows-[auto_1fr] w-full">
      <header className="p-4 py-6 border-b border-dre-200/30 grid grid-cols-[auto_1fr_auto] gap-2">
        <div>
          <Button
            asChild
            className="rounded-full aspect-square p-1 justify-center"
            variant="ghost"
            title="Back to dispute"
          >
            <Link href={`/disputes/${id}`}>
              <ChevronLeftIcon />
            </Link>
          </Button>
        </div>
        <div>
          <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
            Dispute tickets
          </h1>
          <p>Dispute ID: {id}</p>
        </div>
        <CreateTicketDialog asChild dispute={id}>
          <Button className="mt-3">Create ticket</Button>
        </CreateTicketDialog>
      </header>
      <main className="p-4 py-6">
        <ul className="space-y-6">
          {data.tickets.map((ticket) => (
            <Card key={ticket.id} className="p-4 grid grid-cols-[1fr_auto_auto] items-center gap-3">
              <div className="space-y-2">
                <CardTitle>{ticket.subject}</CardTitle>
                <CardDescription>Opened on {ticket.date_created}</CardDescription>
              </div>
              <p>{ticket.status}</p>
              <Button asChild variant="outline">
                <Link href={`./tickets/${ticket.id}`}>Read more...</Link>
              </Button>
            </Card>
          ))}
        </ul>
      </main>
    </div>
  );
}