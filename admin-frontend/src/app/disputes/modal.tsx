import { Button } from "@/components/ui/button";

import { DialogClose, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { X } from "lucide-react";
import Sidebar from "@/components/admin/sidebar";

export default function DisputeDetails({ open = false }: { open?: boolean }) {
  return (
    <Sidebar open={open} className="p-6 md:pl-8 rounded-l-2xl">
      <DialogHeader className="grid grid-cols-[1fr_auto] gap-2 border-b pb-6 border-primary-500/50 space-y-0 items-center">
        <DialogTitle className="p-2">
          Mr. Plastic CC vs. Mr. Plastic & Mining Promotional Goods
        </DialogTitle>
        <div className="flex justify-end items-start">
          <DialogClose asChild>
            <Button variant="ghost" className="rounded-full aspect-square p-2 m-0">
              <X />
            </Button>
          </DialogClose>
        </div>
        <div className="flex gap-2 items-center">
          <Button>Brother what</Button>
          <span>Filed 2 days ago</span>
        </div>

        <p>Case Number: ZA2007-0001</p>
      </DialogHeader>
      <DialogFooter>
        <Button type="submit">Save changes</Button>
      </DialogFooter>
    </Sidebar>
  );
}
