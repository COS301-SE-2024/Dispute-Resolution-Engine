import { getTicketDetails } from "@/lib/api/tickets";

type Props = {
  params: { tid: string };
};

export default async function TicketDetails({ params: { tid } }: Props) {
  const details = await getTicketDetails(parseInt(tid));

  return (
    <main>
      <header>
        <h1 className="text-2xl font-bold">{details.subject}</h1>
      </header>
      <p>{details.body}</p>
    </main>
  );
}
