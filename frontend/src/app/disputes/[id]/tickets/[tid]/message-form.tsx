"use client";

import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { addTicketMessage } from "@/lib/api/tickets";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

const messageSchema = z.object({
  message: z.string().trim().min(1),
});
type MessageData = z.infer<typeof messageSchema>;

export default function MessageForm({ ticket }: { ticket: number }) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<MessageData>({
    resolver: zodResolver(messageSchema),
  });

  async function onSubmit(data: MessageData) {
    await addTicketMessage(ticket, data.message);
  }

  return (
    <Card asChild>
      <form onSubmit={handleSubmit(onSubmit)}>
        <CardHeader>
          <CardTitle>Send a message</CardTitle>
        </CardHeader>
        <CardContent>
          <Textarea {...register("message")} />
        </CardContent>
        <CardFooter className="justify-end">
          {errors.message && (
            <p role="alert" className="text-red-500 grow">
              {errors.message.message}
            </p>
          )}
          <Button>Send</Button>
        </CardFooter>
      </form>
    </Card>
  );
}
