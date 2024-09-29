"use client";
import { Handle, Node, NodeProps, Position, useReactFlow } from "@xyflow/react";
import { FormEvent, ReactNode, useId, useRef, useState } from "react";
import { BookOpenIcon, CirclePlus, CircleX, ClockIcon, Pencil } from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { GraphInstance, TimerDuration, type GraphState } from "@/lib/types";
import EditForm from "./edit-form";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/lib/hooks/use-toast";

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
    let nodes = reactFlow.getNodes();
    let edges = reactFlow.getEdges();
    for (let edge of edges) {
      if (edge.target == data.id) {
        let sourceNode = nodes.find((node) => node.id == edge.source);
        if (sourceNode && sourceNode.data.edges) {
          sourceNode.data.edges = sourceNode.data.edges.filter(
            (handle) => handle.id != edge.sourceHandle
          );
        }
      }
    }
    edges = edges.filter((edge) => edge.target != data.id);
    nodes = nodes.filter((node) => node.id != data.id);
    reactFlow.setNodes(nodes);
    reactFlow.setEdges(edges);
  }

  /** Used to determine when a component the label of a node is being edited */
  const [editing, setEditing] = useState(false);

  function setNodeLabel(value: string) {
    setEditing(false);
    reactFlow.updateNodeData(data.id, {
      label: value,
    });
  }

  function setNodeDescription(value: string) {
    reactFlow.updateNodeData(data.id, {
      description: value,
    });
  }
  function setNodeTimer(dur: TimerDuration) {
    reactFlow.updateNodeData(data.id, {
      timer: dur,
    });
  }

  return (
    <Card className="min-w-48">
      <CardHeader className="p-3 flex gap-1 flex-row items-center">
        {!(data.data.initial ?? false) && (
          <Button
            variant="ghost"
            className="rounded-full p-2 items-center justify-center"
            onClick={deleteNode}
          >
            <CircleX size="1rem" />
          </Button>
        )}
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
      <CardContent style={{ minHeight }} className="relative pt-0 mt-0 flex flex-col items-start">
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
        <DescriptionEditor
          value={data.data.description}
          state={data.data.label}
          asChild
          onValueChange={setNodeDescription}
        >
          <Button variant="ghost" className="text-sm font-normal gap-2">
            <BookOpenIcon size="1rem" />
            Edit description
          </Button>
        </DescriptionEditor>

        <TimerEditor
          state={data.data.label}
          value={data.data.timer}
          onValueChange={setNodeTimer}
          asChild
        >
          <Button variant="ghost" className="text-sm font-normal gap-2">
            <ClockIcon size="1rem" />
            {data.data.timer ? "Edit timer" : "Add timer"}
          </Button>
        </TimerEditor>
        {/* <TimerCheckbox data={data} /> */}
        {/* <EventSection></EventSection> */}
      </CardContent>
    </Card>
  );
}

function DescriptionEditor({
  children,
  asChild,
  state,
  value,
  onValueChange = () => {},
}: {
  children: ReactNode;
  asChild?: boolean;
  state: string;
  value: string;
  onValueChange?: (val: string) => void;
}) {
  const area = useRef<HTMLTextAreaElement>(null);

  function onSubmit() {
    const value = area.current!.value.trim();
    if (value.length === 0) {
      return;
    }
    onValueChange(value);
  }

  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{}Edit description</DialogTitle>
          <DialogDescription>
            Change the description for the &quot;{state}&quot; state
          </DialogDescription>
        </DialogHeader>
        <Textarea ref={area} defaultValue={value} className="resize" />
        <DialogFooter>
          <DialogClose asChild>
            <Button onClick={onSubmit}>Confirm</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function TimerEditor({
  children,
  asChild,
  state,
  value,
  onValueChange = () => {},
}: {
  children: ReactNode;
  asChild?: boolean;
  state: string;
  value?: TimerDuration;
  onValueChange?: (val: TimerDuration) => void;
}) {
  const days = useRef<HTMLInputElement>(null);
  const hours = useRef<HTMLInputElement>(null);
  const minutes = useRef<HTMLInputElement>(null);
  const seconds = useRef<HTMLInputElement>(null);

  const dayId = useId();
  const hourId = useId();
  const minuteId = useId();
  const secondId = useId();

  const { toast } = useToast();
  function onSubmit() {
    try {
      const dayValue = parseInt(days.current!.value.trim());
      const hourValue = parseInt(hours.current!.value.trim());
      const minuteValue = parseInt(minutes.current!.value.trim());
      const secondValue = parseInt(seconds.current!.value.trim());
      onValueChange({
        days: isNaN(dayValue) ? 0 : dayValue,
        hours: isNaN(hourValue) ? 0 : hourValue,
        minutes: isNaN(minuteValue) ? 0 : minuteValue,
        seconds: isNaN(secondValue) ? 0 : secondValue,
      });
    } catch (e: unknown) {
      const error = e as Error;
      toast({
        variant: "error",
        title: "Failed to update timer",
        description: error?.message,
      });
    }
  }

  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{value ? "Edit timer" : "Add timer"}</DialogTitle>
          <DialogDescription>Edit the timer for the &quot;{state}&quot; state</DialogDescription>
        </DialogHeader>
        <div className="grid grid-cols-[auto_1fr] items-center gap-x-2 gap-y-4">
          <Label htmlFor={dayId}>Days</Label>
          <Input defaultValue={value?.days} type="number" id={dayId} ref={days} name="days" />

          <Label htmlFor={hourId}>Hours</Label>
          <Input defaultValue={value?.hours} type="number" id={hourId} ref={hours} name="hours" />

          <Label htmlFor={minuteId}>Minutes</Label>
          <Input
            defaultValue={value?.minutes}
            type="number"
            id={minuteId}
            ref={minutes}
            name="minutes"
          />

          <Label htmlFor={secondId}>Second</Label>
          <Input
            defaultValue={value?.seconds}
            type="number"
            id={secondId}
            name="seconds"
            ref={seconds}
          />
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button onClick={onSubmit}>Confirm</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
