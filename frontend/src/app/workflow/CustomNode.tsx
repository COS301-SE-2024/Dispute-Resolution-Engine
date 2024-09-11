import { Handle, Node, NodeProps, Position } from "@xyflow/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import TimerCheckbox from "./TimerCheckbox";
import EventSection from "./EventSection";
import { ReactNode, useEffect, useId, useState } from "react";
import { eventType } from "@/lib/types";

const events = [
  {id: "a"},
  {id: "b"},
  {id: "c"},
  // {id: "d"},
  // {id: "e"},
  // {id: "f"}, 
  // {id: "g"},
  // {id: "h"},
  // {id: "i"},
]
export type CustomNodeType = Node<{
edges? : any
label? : any
}, "customNode">
export default function CustomNode(data: NodeProps<CustomNodeType>) {
  // console.log(data)
  console.log("rerender")
  const numHandles = events.length
  const gap = 30
  const handles = events.map((event, index) => {
    return <Handle type="target" key={index} id={event.id} style={{height: 20, width: 20, color: "blue  ",top: (140 - numHandles * gap/4) + (index * gap)}} position={Position.Right} />
  })
  return (
    <div className="bg-opacity-100">
      {handles}
      {/* <Handle type="target" id="a" position={Position.Right} />
      <Handle type="target" id="b" style={handleStyle} position={Position.Right} /> */}
      <Card className="dark:bg-black min-w-48">
        <CardHeader className="p-3 text-center">
          {/* TODO: USE state to actually change the label */}
          <CardTitle contentEditable="true" className="text-3xl" suppressContentEditableWarning={true}>
            {data.data.label as ReactNode}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <TimerCheckbox data={data} />
          <EventSection></EventSection>
        </CardContent>
      </Card>
      <Handle type="source" position={Position.Left} id="a" />
    </div>
  );
}
