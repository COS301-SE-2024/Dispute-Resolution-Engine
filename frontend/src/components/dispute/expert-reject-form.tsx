"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "../ui/button";
import { Form, FormMessage, FormField, FormSubmit } from "../ui/form-server";
import { ExpertRejectData } from "@/lib/schema/dispute";
import { rejectExpert } from "@/lib/actions/dispute";
import { Textarea } from "../ui/textarea";
import { ReactNode, useId } from "react";
import { DialogDescription } from "@radix-ui/react-dialog";

const RejectForm = Form<ExpertRejectData>;
const RejectMessage = FormMessage<ExpertRejectData>;
const RejectField = FormField<ExpertRejectData>;

export default function ExpertRejectForm({
  name,
  expertId,
  disputeId,
  children,
  asChild,
}: {
  name?: string;
  expertId: string;
  disputeId: string;
  children: ReactNode;
  asChild?: boolean;
}) {
  const reasonId = useId();

  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Object to Expert Assignment</DialogTitle>
          {name && <DialogDescription>Objecting to {name}</DialogDescription>}
        </DialogHeader>
        <RejectForm action={rejectExpert} className="space-y-2 w-full">
          <input type="hidden" name="expert_id" value={expertId} />
          <input type="hidden" name="dispute_id" value={disputeId} />
          <RejectField id={reasonId} name="reason" label="Reason" className="col-span-2">
            {name ? (
              <Textarea
                id={reasonId}
                placeholder={`Why do you object to ${name}? (min. 20 characters)`}
                name="reason"
              />
            ) : (
              <Textarea
                id={reasonId}
                placeholder={"Why do you object to this assignment? (min. 20 characters)"}
                name="reason"
              />
            )}
          </RejectField>
          <div className="flex justify-end gap-2 items-center">
            <RejectMessage />
            <FormSubmit>Submit</FormSubmit>
          </div>
        </RejectForm>
      </DialogContent>
    </Dialog>
  );
}
