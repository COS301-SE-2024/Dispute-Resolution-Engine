"use client";

import { Expert } from "@/lib/interfaces/dispute";
import { Badge } from "@/components/ui/badge";
import ExpertRejectForm from "./expert-reject-form";
import { Button } from "../ui/button";

export interface ExpertItemProps extends Expert {
  dispute_id: string;
}

export default function ExpertItem(props: ExpertItemProps) {
  return (
    <section className="p-2 rounded-lg">
      <div className="flex items-center gap-2">
        <h4 className="text-xl font-bold">{props.full_name}</h4>
        <Badge>{props.role}</Badge>
      </div>
      <div className="flex justify-between flex-wrap gap-3">
        <dl className="grid grid-cols-2">
          <dt className="font-semibold">Email</dt>
          <dd>{props.email}</dd>

          <dt className="font-semibold">Phone</dt>
          <dd>{props.phone}</dd>
        </dl>
        <div className="flex gap-2">
          <ExpertRejectForm
            disputeId={props.dispute_id}
            expertId={props.id}
            name={props.full_name}
            asChild
          >
            <Button variant="destructive">Object to expert</Button>
          </ExpertRejectForm>
        </div>
      </div>
    </section>
  );
}
