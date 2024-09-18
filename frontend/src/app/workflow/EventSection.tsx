import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { CircleMinus, CirclePlus } from "lucide-react";

type DisputeEvent = {
  name: string;
  trigger: string;
};
const events: DisputeEvent[] = [
  {
    name: "First",
    trigger: "Submission received",
  },
  {
    name: "First",
    trigger: "Timer",
  },
];
export default function EventSection() {
  return (
    <div className="grid grid-cols-[auto_1fr] gap-1 grid-rows-2">
      <div className="row-span-2 flex items-center justify-center">
        <Button>
          <CircleMinus />
        </Button>
      </div>
      <div className="flex justify-center items-center">
        <h3 className="text-2xl">First</h3>
      </div>
      <div className="flex items-center pl-4">
        <Input placeholder="Submission received"></Input>
      </div>
      <Button>
        <CirclePlus/>
      </Button>
    </div>
  );
}
