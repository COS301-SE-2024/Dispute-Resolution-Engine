"use client";
import { Handle, Node, NodeProps, Position, useReactFlow } from "@xyflow/react";
import { FormEvent, ReactNode, useId, useRef, useState } from "react";
import {
  BookOpenIcon,
  CirclePlus,
  CircleX,
  ClockIcon,
  Pencil,
  PlusIcon,
  TrashIcon,
} from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { GraphInstance, TimerDuration, type GraphState } from "@/lib/types";
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
import EditForm from "./edit-form";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

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
  function setNodeTimer(dur: TimerDuration | undefined) {
    reactFlow.updateNodeData(data.id, {
      timer: dur,
    });
  }
  function removeNodeTimer() {
    reactFlow.updateNodeData(data.id, {
      timer: undefined,
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
          className="flex items-center justify-center text-center"
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

        <div className="flex">
          <TimerEditor
            state={data.data.label}
            value={data.data.timer}
            onValueChange={setNodeTimer}
            asChild
          >
            <Button variant="ghost" className="text-sm font-normal gap-2 grow">
              {data.data.timer ? (
                <>
                  <ClockIcon size="1rem" />
                  Edit timer
                </>
              ) : (
                <>
                  <PlusIcon size="1rem" />
                  Add timer
                </>
              )}
            </Button>
          </TimerEditor>

          {data.data.timer && (
            <Button
              variant="ghost"
              className="text-sm font-normal gap-2 grow"
              title="Remove timer"
              onClick={() => removeNodeTimer()}
            >
              <TrashIcon size="1rem" className="text-red-500" />
            </Button>
          )}
        </div>
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

const timerSchema = z.object({
  days: z.coerce
    .number({
      message: "Invalid day",
    })
    .optional(),
  hours: z.coerce
    .number({
      message: "Invalid hours",
    })
    .optional(),
  minutes: z.coerce
    .number({
      message: "Invalid minutes",
    })
    .optional(),
  seconds: z.coerce
    .number({
      message: "Invalid seconds",
    })
    .optional(),
});
type TimerData = z.infer<typeof timerSchema>;

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
  onValueChange?: (val: TimerDuration | undefined) => void;
}) {
  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<TimerData>({
    resolver: zodResolver(timerSchema),
    defaultValues: {
      days: value?.days ?? 0,
      hours: value?.hours ?? 0,
      minutes: value?.minutes ?? 0,
      seconds: value?.seconds ?? 0,
    },
  });

  const formId = useId();
  const dayId = useId();
  const hourId = useId();
  const minuteId = useId();
  const secondId = useId();

  const [open, setOpen] = useState(false);
  function onSubmit(data: TimerData) {
    const value = {
      days: data.days ?? 0,
      hours: data.hours ?? 0,
      minutes: data.minutes ?? 0,
      seconds: data.seconds ?? 0,
    };

    const total = value.days + value.hours + value.minutes + value.seconds;
    if (total === 0) {
      onValueChange(undefined);
    } else {
      onValueChange(value);
    }
    setOpen(false);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{value ? "Edit timer" : "Add timer"}</DialogTitle>
          <DialogDescription>Edit the timer for the &quot;{state}&quot; state</DialogDescription>
        </DialogHeader>
        <form
          id={formId}
          onSubmit={handleSubmit(onSubmit)}
          className="grid grid-cols-2 items-center gap-x-2 gap-y-4"
        >
          <div>
            <Label htmlFor={dayId}>Days</Label>
            <Input id={dayId} {...register("days")} />
            {errors.days && (
              <p role="alert" className="text-red-500 text-sm">
                {errors.days.message}
              </p>
            )}
          </div>

          <div>
            <Label htmlFor={hourId}>Hours</Label>
            <Input id={hourId} {...register("hours")} />
            {errors.hours && (
              <p role="alert" className="text-red-500 text-sm">
                {errors.hours.message}
              </p>
            )}
          </div>

          <div>
            <Label htmlFor={minuteId}>Minutes</Label>
            <Input id={minuteId} {...register("minutes")} />
            {errors.minutes && (
              <p role="alert" className="text-red-500 text-sm">
                {errors.minutes.message}
              </p>
            )}
          </div>

          <div>
            <Label htmlFor={secondId}>Seconds</Label>
            <Input id={secondId} {...register("seconds")} />
            {errors.seconds && (
              <p role="alert" className="text-red-500 text-sm">
                {errors.seconds.message}
              </p>
            )}
          </div>
        </form>
        <DialogFooter>
          <Button form={formId}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
