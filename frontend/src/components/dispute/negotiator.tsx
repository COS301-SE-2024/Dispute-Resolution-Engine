"use client";

import { useId } from "react";

import { Expert } from "@/lib/interfaces/dispute";

import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Form, FormField, FormMessage, FormSubmit } from "../ui/form-server";
import { ExpertRejectData } from "@/lib/schema/dispute";
import { approveExpert, rejectExpert } from "@/lib/actions/dispute";

const RejectForm = Form<ExpertRejectData>;
const RejectMessage = FormMessage<ExpertRejectData>;
const RejectField = FormField<ExpertRejectData>;

const ApproveForm = Form<ExpertRejectData>;

export interface ExpertItemProps extends Expert {
    dispute_id: string;
}

export default function ExpertItem(props: ExpertItemProps) {
    const areaId = useId();

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
                    <ApproveForm action={approveExpert}>
                        <input type="hidden" name="dispute_id" value={props.dispute_id} />
                        <input type="hidden" name="expert_id" value={props.id} />
                        <Button type="submit">Approve</Button>
                    </ApproveForm>
                    <Dialog>
                        <DialogTrigger asChild>
                            <Button variant="destructive">Reject</Button>
                        </DialogTrigger>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Reject {props.full_name}</DialogTitle>
                            </DialogHeader>
                            <RejectForm action={rejectExpert} className="space-y-2 w-full">
                                <input type="hidden" name="expert_id" value={props.id} />
                                <input type="hidden" name="dispute_id" value={props.dispute_id} />
                                <RejectField
                                    id={areaId}
                                    name="reason"
                                    label="Reason"
                                    className="col-span-2"
                                >
                                    <Textarea id={areaId} placeholder={`Why do you object to ${props.full_name}? (min. 20 characters)`} name="reason" />
                                </RejectField>
                                <div className="flex justify-end gap-2 items-center">
                                    <RejectMessage />
                                    <FormSubmit>Reject</FormSubmit>
                                </div>
                            </RejectForm>
                        </DialogContent>
                    </Dialog>
                </div>
            </div>
        </section>
    );
}
