"use client";
import { Handle, Node, NodeProps, Position, useReactFlow } from "@xyflow/react";
import { FormEvent, ReactNode, useRef, useState } from "react";
import { CirclePlus, CircleX, Pencil } from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { GraphInstance, type GraphState } from "@/lib/types";
import EditForm from "./edit-form";

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

export default function CustomNode(data: NodeProps<GraphState>) {
  const events = data.data.edges;
  const numHandles = events.length;

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
      ></Handle>
    );
  });

  const reactFlow: GraphInstance = useReactFlow();

  function deleteNode() {
    reactFlow.setNodes((nodes) => nodes.filter((node) => node.id !== data.id));
  }

  /** Used to determine when a component the label of a node is being edited */
  const [editing, setEditing] = useState(false);

  function setNodeLabel(value: string) {
    setEditing(false);
    reactFlow.updateNodeData(data.id, {
      label: value,
    });
  }

  return (
    <Card className="min-w-48">
      <CardHeader className="p-3 flex gap-1 flex-row items-center">
        <Button
          variant="ghost"
          className="rounded-full p-2 items-center justify-center"
          onClick={deleteNode}
        >
          <CircleX size="1rem" />
        </Button>
        {editing ? (
          <EditForm
            value={data.data.label}
            onCommit={setNodeLabel}
            onCancel={() => setEditing(false)}
          />
        ) : (
          <>
            <CardTitle className="grow text-base">{data.data.label as ReactNode}</CardTitle>
            <Button variant="ghost" className="rounded-full p-2" onClick={() => setEditing(true)}>
              <Pencil size="1rem" />
            </Button>
          </>
        )}
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
          +
        </Handle>
        <Handle type="target" position={Position.Left} id="a" style={handleStyle} />
        {/* <TimerCheckbox data={data} /> */}
        {/* <EventSection></EventSection> */}
      </CardContent>
    </Card>
  );
}
