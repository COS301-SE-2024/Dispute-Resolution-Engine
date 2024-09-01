import { Button } from "@/components/ui/button";

import {
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Download, EllipsisVertical, FileText, Trash, X } from "lucide-react";
import Sidebar from "@/components/admin/sidebar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

function Evidence() {
  return (
    <li className="grid grid-cols-[auto_1fr_auto] gap-2 items-center px-3 py-2 border border-primary-500/30 rounded-md">
      <FileText className="stroke-primary-500" size="1.7rem" />
      <div>
        <span className="truncate">Evidence name</span> <br />
        <span className="truncate text-white/50">Submission date</span>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="rounded-full p-2">
            <EllipsisVertical />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem>
            <Download className="mr-2" />
            <span>Download file</span>
          </DropdownMenuItem>
          <DropdownMenuItem className="text-red-500">
            <Trash className="mr-2" />
            <span>Delete evidence</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </li>
  );
}

export default function DisputeDetails({ open = false }: { open?: boolean }) {
  return (
    <Sidebar open={open} className="p-6 md:pl-8 rounded-l-2xl space-y-6">
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
      <Card>
        <CardHeader>
          <CardTitle>Overview</CardTitle>
          <CardDescription>Description here</CardDescription>
        </CardHeader>
        <CardContent>
          <strong>Evidence</strong>
          <ul className="grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))] gap-2">
            <Evidence />
            <Evidence />
            <Evidence />
            <Evidence />
            <Evidence />
            <Evidence />
          </ul>
        </CardContent>
      </Card>
    </Sidebar>
  );
}
