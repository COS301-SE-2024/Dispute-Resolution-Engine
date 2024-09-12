"use client"
import { Handle, Node, NodeProps, Position, ReactFlowInstance, useReactFlow } from "@xyflow/react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import TimerCheckbox from "./TimerCheckbox";
import EventSection from "./EventSection";
import { ReactNode, useEffect, useId, useRef, useState } from "react";
import { eventType } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { CircleX } from "lucide-react";

import flow from "@xyflow/react";

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

/** The diameter (in pixels) of a single handle */
const handleDiameter = 20;

/** The gap (in pixels) between handles */
const handleGap = 20;

/** Default styles applied to every handle */
const handleStyle = {
  height: handleDiameter,
  width: handleDiameter,
};

/** Calculates a handle's offset given it's index */
const offset = (i: number) => i * (handleDiameter + handleGap);

export default function CustomNode(data: NodeProps<CustomNodeType>) {
  const events = data.data.edges;
  const numHandles = events.length;
  console.log("rerender", numHandles);

  const minHeight = offset(numHandles + 1);

  const handles = events.map((event, index) => {
    return (
      <Handle
        type="source"
        key={index}
        id={event.id}
        style={{
          ...handleStyle,
          color: "blue  ",
          top: offset(index),
        }}
        position={Position.Right}
      >
        {event.id}
      </Handle>
    );
  });
  const flowInst = useReactFlow()
  const deleteNode = function(){
    const nodes = flowInst.getNodes()
    flowInst.setNodes((nodes) => nodes.filter((node) => node.id !== data.id));
  }
  return (
    <Card className="min-w-48">
      {/* <Handle type="target" id="a" position={Position.Right} />
      <Handle type="target" id="b" style={handleStyle} position={Position.Right} /> */}
      <CardHeader className="p-3 text-center">
        {/* TODO: USE state to actually change the label */}
      <button className="nodrag nopan" onClick={deleteNode}>
        <CircleX></CircleX>
      </button>
        <CardTitle className="text-3xl" suppressContentEditableWarning={true}>
          {data.data.label as ReactNode}
        </CardTitle>
      </CardHeader>
      <CardContent style={{ minHeight }} className="relative">
        {handles}
        <Handle
          type="source"
          id="new"
          style={{
            ...handleStyle,
            color: "white",
            top: offset(numHandles),
          }}
          position={Position.Right}
        >
          Event_Name
        </Handle>
        <Handle type="target" position={Position.Left} id="a" style={handleStyle} />
        {/* <TimerCheckbox data={data} /> */}
        {/* <EventSection></EventSection> */}
      </CardContent>
    </Card>
  );
}
