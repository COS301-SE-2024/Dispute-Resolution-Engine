"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Form, FormMessage, FormField, FormSubmit } from "../ui/form-server";
import { DisputeDecisionData } from "@/lib/schema/dispute";
import { rejectExpert, uploadDecision } from "@/lib/actions/dispute";
import { ReactNode, useId } from "react";
import { DialogDescription } from "@radix-ui/react-dialog";
import { Input } from "../ui/input";
import { SelectProps } from "@radix-ui/react-select";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

import { DISPUTE_DECISION } from "@/lib/interfaces/dispute";

const DecisionForm = Form<DisputeDecisionData>;
const DecisionMessage = FormMessage<DisputeDecisionData>;
const DecisionField = FormField<DisputeDecisionData>;

export default function DisputeDecisionForm({
  disputeId,
  children,
  asChild,
}: {
  name?: string;
  disputeId: string;
  children: ReactNode;
  asChild?: boolean;
}) {
  const decisionId = useId();
  const writeupId = useId();

  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Render a decision for a dispute</DialogTitle>
          <DialogDescription>This will close a dispute, based on your decision</DialogDescription>
        </DialogHeader>
        <DecisionForm action={uploadDecision} className="space-y-2 w-full">
          <input type="hidden" name="dispute_id" value={disputeId} />
          <DecisionField id={decisionId} name="decision" label="Decision" className="col-span-2">
            <DecisionSelect id={decisionId} name="decision" />
          </DecisionField>
          <DecisionField id={writeupId} name="writeup" label="Writeup" className="col-span-2">
            <Input type="file" id={writeupId} name="writeup" />
          </DecisionField>
          <div className="flex justify-end gap-2 items-center">
            <DecisionMessage />
            <FormSubmit>Submit</FormSubmit>
          </div>
        </DecisionForm>
      </DialogContent>
    </Dialog>
  );
}

function DecisionSelect({
  id,
  ...props
}: SelectProps & {
  id: string;
}) {
  return (
    <Select {...props}>
      <SelectTrigger>
        <SelectValue id={id} placeholder="Select a decision" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {DISPUTE_DECISION.map((gen) => (
            <SelectItem key={gen} value={gen}>
              {gen}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
