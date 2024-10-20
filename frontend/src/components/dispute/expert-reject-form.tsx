"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "../ui/button";
import { FormMessage, FormField } from "../ui/form-client";
import { ExpertRejectData, expertRejectSchema } from "@/lib/schema/dispute";
import { rejectExpert } from "@/lib/actions/dispute";
import { Textarea } from "../ui/textarea";
import { ReactNode, useId } from "react";
import { DialogDescription } from "@radix-ui/react-dialog";
import { FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";

export const RejectMessage = FormMessage<ExpertRejectData>;
export const RejectField = FormField<ExpertRejectData>;

export default function ExpertRejectForm({
  expertName,
  expertId,
  disputeId,
  children,
  asChild,
}: {
  expertName: string;
  expertId: string;
  disputeId: string;
  children: ReactNode;
  asChild?: boolean;
}) {
  const form = useForm<ExpertRejectData>({
    resolver: zodResolver(expertRejectSchema),
  });
  const { register, handleSubmit, setError } = form;
  const router = useRouter();

  async function onSubmit(data: ExpertRejectData) {
    const res = await rejectExpert(data)
    if (res.data) {
      router.push(`/disputes/${disputeId}/tickets/${res.data}`);  
    }
    // rejectExpert(data)
    //   .then((id) => {
    //     router.push(`/disputes/${disputeId}/tickets/${id}`);
    //   })
    //   .catch((e: Error) => {
    //     setError("root", {
    //       type: "custom",
    //       message: e.message,
    //     });
    //   });
  }

  const formId = useId();
  const reasonId = useId();

  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Object to Expert Assignment</DialogTitle>
          <DialogDescription>Objecting to {expertName}</DialogDescription>
        </DialogHeader>

        <FormProvider {...form}>
          <div className="grid grid-rows-[1fr_auto] gap-5">
            <form id={formId} onSubmit={handleSubmit(onSubmit)}>
              <input type="hidden" value={disputeId} {...register("dispute_id")} />
              <input type="hidden" value={expertId} {...register("expert_id")} />

              <RejectField name="reason" label="Reason" className="col-span-2" id={reasonId}>
                <Textarea
                  id={reasonId}
                  placeholder={`Tell us what's on your mind!`}
                  {...register("reason")}
                />
              </RejectField>
            </form>
            <div className="flex justify-end gap-2 items-center">
              <RejectMessage />
              <Button form={formId} type="submit" variant="outline">
                Submit objection
              </Button>
            </div>
          </div>
        </FormProvider>
      </DialogContent>
    </Dialog>
  );
}
