import { XIcon } from "lucide-react";
import { Button } from "../ui/button";
import { DialogClose, DialogHeader, DialogTitle } from "../ui/dialog";
import { HTMLAttributes } from "react";

export interface Props extends HTMLAttributes<HTMLDivElement> {
  title: string;
}

export default function SidebarHeader({ title, ...props }: Props) {
  return (
    <DialogHeader className="grid grid-cols-[1fr_auto] gap-2 border-b pb-6 mb-6 border-primary-500/50 space-y-0 items-center">
      <DialogTitle className="p-2">{title}</DialogTitle>
      <div className="flex justify-end items-start">
        <DialogClose asChild>
          <Button variant="ghost" className="rounded-full aspect-square p-2 m-0">
            <XIcon />
          </Button>
        </DialogClose>
      </div>
      <div {...props} />
    </DialogHeader>
  );
}
