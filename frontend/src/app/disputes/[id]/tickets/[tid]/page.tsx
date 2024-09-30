import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { getTicketDetails } from "@/lib/api/tickets";
import MessageForm from "./message-form";

type Props = {
  params: { tid: string };
};

export default async function TicketDetails({ params: { tid } }: Props) {
  const details = await getTicketDetails(parseInt(tid));

  return (
    <div className="grid grid-rows-[auto_1fr] w-full">
      <header className="p-4 py-6 border-b border-dre-200/30">
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">
          {details.subject}
        </h1>
        <p>Status: {details.status}</p>
      </header>
      <main className="p-4 space-y-4 overflow-y-auto">
        <Card>
          <CardHeader>
            <CardTitle>Description</CardTitle>
            <CardDescription>{details.body}</CardDescription>
          </CardHeader>
        </Card>
        {details.messages.map((ticket) => (
          <Card key={ticket.id}>
            <CardHeader>
              <CardTitle>{ticket.user.full_name}</CardTitle>
              <CardDescription>Sent on {ticket.date_sent}</CardDescription>
            </CardHeader>
            <CardContent asChild>
              <p>{ticket.message}</p>
            </CardContent>
          </Card>
        ))}
        <MessageForm ticket={parseInt(tid)} />
      </main>
    </div>
  );
}
