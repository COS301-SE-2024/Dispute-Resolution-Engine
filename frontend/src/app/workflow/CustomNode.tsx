import { useCallback, useState } from "react";
import { Handle, Position } from "@xyflow/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { disputeDuration } from "@/lib/types";

const handleStyle = { left: 10 };


export default function CustomNode(data: any) {
  // console.log(data)
  const [duration, setDuration] = useState<disputeDuration>({days: 0, hours: 10, minutes: 20, seconds: 30})
  data.data.duration = duration;
  const handleInputChange = (evt: React.ChangeEvent<HTMLInputElement>, unit: keyof disputeDuration) => {
    const value = parseInt(evt.target.value, 10) || 0;
    setDuration((prevDuration) => ({
      ...prevDuration,
      [unit]: value,
    }));
  };
  const arr = ["days", "hours", "minutes", "seconds"];
  const inputs = arr.map((currInput) => {
    return (
      <div key={currInput}>
      <Label htmlFor={currInput + "Inp"}>
        {currInput.charAt(0).toUpperCase() + currInput.slice(1)}
      </Label>
      <Input
        id={currInput + "Inp"}
        type="number"
        className="h-2 w-16 p-3"
        value={duration[currInput as keyof disputeDuration]}
        onChange={(evt) => handleInputChange(evt, currInput as keyof disputeDuration)}
      />
    </div>
    );
  });
  return (
    <div className="bg-opacity-100">
      <Handle type="target" position={Position.Top} />
      <Card className="dark:bg-black min-w-48">
        <CardHeader className="p-3 text-center">
          {/* TODO: USE state to actually change the label */}
          <CardTitle contentEditable="true" suppressContentEditableWarning={true}>{data.data.label}</CardTitle>
        </CardHeader>
        <CardContent>
          {/* <div className="grid grid-cols-2 items-center gap-3">{inputs}</div> */}
          <Accordion type="single" collapsible className="w-full">
            <AccordionItem value="item-1">
              <AccordionTrigger className="pt-0">Timer: {duration.days == 0 ? "" : duration.days.toString()+ "d"} {duration.hours == 0 ? "" : duration.hours.toString() + "h"} {duration.minutes == 0 ?  "" : duration.minutes.toString() + "m"} {duration.seconds == 0 ? "" : duration.seconds.toString() + "s"}</AccordionTrigger>
              <AccordionContent className="grid grid-cols-2 items-center gap-3">{ inputs }</AccordionContent>
            </AccordionItem>
          </Accordion>
        </CardContent>
      </Card>
      <Handle type="source" position={Position.Bottom} id="a" />
    </div>
  );
}
