import { getTicketSummaries } from "@/lib/api/tickets";
import Link from "next/link";

type Props = {
  params: { id: string };
};

export default async function TicketsPage({ params: { id } }: Props) {
  const data = await getTicketSummaries(parseInt(id));
  return (
    <main>
      <ul>
        {data.tickets.map((ticket) => (
          <li key={ticket.id}>
            <Link href={`./tickets/${ticket.id}`}>{ticket.subject}</Link>
          </li>
        ))}
      </ul>
    </main>
  );
}
