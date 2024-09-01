"use client";

import { Dialog, DialogPortal, DialogContent } from "@radix-ui/react-dialog";
import { DialogOverlay } from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";
import { useRouter } from "next/navigation";

export default function Sidebar({
  open = false,
  children,
  className,
}: {
  open?: boolean;
  children?: ReactNode;
  className?: string;
}) {
  const router = useRouter();
  return (
    <Dialog
      open={open}
      onOpenChange={(open) => {
        if (!open) router.back();
      }}
    >
      <DialogPortal>
        <DialogOverlay />
        <DialogContent
          className={cn(
            "fixed z-50 bottom-0 right-0 w-full h-[90%] md:w-[60%] md:h-[100%]",
            "bg-dre-bg-light dark:bg-dre-bg-dark",
            "data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
            "data-[state=closed]:slide-out-to-right-1/2 data-[state=open]:slide-in-from-right-1/2",
            className
          )}
        >
          {children}
        </DialogContent>
      </DialogPortal>
    </Dialog>
  );
}
