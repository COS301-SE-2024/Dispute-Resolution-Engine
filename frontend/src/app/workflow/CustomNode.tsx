import { useCallback } from "react";
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

const handleStyle = { left: 10 };

export default function CustomNode(data: any) {
  const onChange = useCallback((evt: any) => {
    console.log(evt.target.value);
  }, []);
  const arr = ["days", "hours", "minutes", "seconds"];
  const inputs = arr.map((currInput) => {
    return (
      <>
        <Label htmlFor={currInput + "Inp"}>
          {currInput.charAt(0).toUpperCase() + currInput.slice(1)}
        </Label>
        <Input id={currInput + "Inp"} type="number" className="h-2 w-16 p-3"></Input>
      </>
    );
  });
  return (
    <div className="bg-opacity-100">
      <Handle type="target" position={Position.Top} />
      <Card className="dark:bg-black min-w-48">
        <CardHeader className="p-3 text-center">
          <CardTitle contentEditable="true">Node A</CardTitle>
        </CardHeader>
        <CardContent>
          {/* <div className="grid grid-cols-2 items-center gap-3">{inputs}</div> */}
          <Accordion type="single" collapsible className="w-full">
            <AccordionItem value="item-1">
              <AccordionTrigger className="pt-0">Timer: 4h3m</AccordionTrigger>
              <AccordionContent className="grid grid-cols-2 items-center gap-3">{ inputs }</AccordionContent>
            </AccordionItem>
          </Accordion>
        </CardContent>
      </Card>
      <Handle type="source" position={Position.Bottom} id="a" />
    </div>
  );
}
