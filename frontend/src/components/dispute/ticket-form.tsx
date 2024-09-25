import { CreateTicketData } from "@/lib/schema/ticket";
import { Form, FormField, FormMessage, FormSubmit } from "../ui/form-server";
import { ReactNode, useId } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Textarea } from "../ui/textarea";
import { createTicket } from "@/lib/actions/ticket";
import { Input } from "../ui/input";

const TicketForm = Form<CreateTicketData>;
const TicketMessage = FormMessage<CreateTicketData>;
const TicketField = FormField<CreateTicketData>;

export default function CreateTicketDialog({
  dispute,
  children,
  asChild,
}: {
  dispute: string;
  children: ReactNode;
  asChild?: boolean;
}) {
  const subjectId = useId();
  const bodyId = useId();

  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create new ticket</DialogTitle>
          <DialogDescription>Got something on your mind? Tell us!</DialogDescription>
        </DialogHeader>
        <TicketForm action={createTicket} className="space-y-2 w-full">
          <input type="hidden" name="dispute" value={dispute} />
          <TicketField id={subjectId} name="subject" label="Subject">
            <Input
              id={subjectId}
              name="subject"
              placeholder="Why are you creating a ticket? (required)"
            />
          </TicketField>
          <TicketField id={bodyId} name="subject" label="Subject">
            <Textarea
              id={bodyId}
              placeholder="Provide some more details (min. characters)"
              name="body"
            />
          </TicketField>
          <div className="flex justify-end gap-2 items-center">
            <TicketMessage />
            <FormSubmit>Submit</FormSubmit>
          </div>
        </TicketForm>
      </DialogContent>
    </Dialog>
  );
}
