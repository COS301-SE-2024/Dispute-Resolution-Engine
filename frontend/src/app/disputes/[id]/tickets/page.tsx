import DisputeHeader from "@/components/dispute/dispute-header";
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
import Link from "next/link";

type Props = {
  params: { id: string };
};

export default async function TicketsPage({ params: { id } }: Props) {
  const data = await getTicketSummaries(parseInt(id));
  return (
    <div className="grid grid-rows-[auto_1fr] w-full">
      <header className="p-4 py-6 border-b border-dre-200/30">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
          Dispute tickets
        </h1>
        <p>Dispute ID: {id}</p>
      </header>
      <main className="p-4 py-6">
        <ul className="space-y-6">
          {data.tickets.map((ticket) => (
            <Card key={ticket.id} className="p-4 grid grid-cols-[1fr_auto] items-center">
              <div className="space-y-2">
                <CardTitle>{ticket.subject}</CardTitle>
                <CardDescription>Opened on {ticket.date_created}</CardDescription>
              </div>
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
