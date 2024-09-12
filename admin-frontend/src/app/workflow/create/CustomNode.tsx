import { Handle, Node, NodeProps, Position } from "@xyflow/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import TimerCheckbox from "./TimerCheckbox";
import EventSection from "./EventSection";
import { ReactNode, useEffect, useId, useState } from "react";
import { eventType } from "@/lib/types";

// const events = [
//   {id: "a"},
//   {id: "b"},
//   {id: "c"},
//   // {id: "d"},
//   // {id: "e"},
//   // {id: "f"},
//   // {id: "g"},
//   // {id: "h"},
//   // {id: "i"},
// ]
export type CustomNodeType = Node<
  {
    edges: [{ id: string }];
    label?: any;
  },
  "customNode"
>;

export default function CustomNode(data: NodeProps<CustomNodeType>) {
  // console.log(data)
  const fullHeigh = 280;
  const events = data.data.edges;
  const numHandles = events.length;
  console.log("rerender", numHandles);
  const gap = 30;
  const handles = events.map((event, index) => {
    return (
      <Handle
        type="source"
        key={index}
        id={event.id}
        style={{
          height: 10,
          width: 10,
          color: "blue  ",
          top: 40 - (numHandles * gap) / 4 + index * gap,
        }}
        position={Position.Right}
      />
    );
  });

  return (
    <div className="bg-opacity-100">
      {handles}
      <Handle
        type="source"
        id="new"
        style={{
          height: 20,
          width: 20,
          color: "blue  ",
          top: 40 - (numHandles * gap) / 4 + numHandles * gap,
        }}
        position={Position.Right}
      />
      {/* <Handle type="target" id="a" position={Position.Right} />
      <Handle type="target" id="b" style={handleStyle} position={Position.Right} /> */}
      <Card className="min-w-48">
        <CardHeader className="p-3 text-center">
          {/* TODO: USE state to actually change the label */}
          <CardTitle
            contentEditable="true"
            className="text-3xl"
            suppressContentEditableWarning={true}
          >
            {data.data.label as ReactNode}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {/* <TimerCheckbox data={data} /> */}
          {/* <EventSection></EventSection> */}
        </CardContent>
      </Card>
      <Handle type="target" position={Position.Left} id="a" />
    </div>
  );
}
