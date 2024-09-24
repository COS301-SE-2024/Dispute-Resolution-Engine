import PageHeader from "@/components/admin/page-header";
import { Card, CardContent } from "@/components/ui/card";
import { TicketTable } from "./table";
import { z } from "zod";
import TicketDetails from "./details";
import { Ticket } from "@/lib/types/tickets";

const searchSchema = z.object({
  id: z.string().optional(),
});

const ticket = {
  id: "0",
  user: { id: "0", full_name: "John Doe" },
  date_created: "2023-09-01T10:30:00Z",
  subject: "Unable to access account",
  status: "Open",
  body: "I cannot log in to my account. It keeps saying incorrect password.",
  messages: [
    {
      id: "0",
      user: { id: "0", full_name: "John Doe" },
      date_sent: "2023-09-01T10:31:00Z",
      message: "I cannot log in to my account. It keeps saying incorrect password.",
    },
  ],
} satisfies Ticket;

export default function Tickets({ searchParams }: { searchParams: unknown }) {
  const { data: params, error: searchError } = searchSchema.safeParse(searchParams);
  if (!params) {
    throw new Error(JSON.stringify(searchError));
  }

  return (
    <>
      {params.id && <TicketDetails details={ticket} />}
      <div className="flex flex-col">
        <PageHeader label="Tickets" />
        <Card>
          <CardContent>
            <TicketTable />
          </CardContent>
        </Card>
      </div>
    </>
  );
}
