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
import TimerCheckbox from "./TimerCheckbox";
import EventSection from "./EventSection";

const handleStyle = { left: 10 };

export default function CustomNode(data: any) {
  // console.log(data)
  return (
    <div className="bg-opacity-100">
      <Handle type="target" position={Position.Top} />
      <Card className="dark:bg-black min-w-48">
        <CardHeader className="p-3 text-center">
          {/* TODO: USE state to actually change the label */}
          <CardTitle contentEditable="true" className="text-3xl" suppressContentEditableWarning={true}>
            {data.data.label}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <TimerCheckbox data={data} />
          <EventSection></EventSection>
        </CardContent>
      </Card>
      <Handle type="source" position={Position.Bottom} id="a" />
    </div>
  );
}
