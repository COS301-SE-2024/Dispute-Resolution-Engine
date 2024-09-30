import {
  BaseEdge,
  Edge,
  EdgeProps,
  EdgeLabelRenderer,
  Position,
  getSmoothStepPath,
  useReactFlow,
  useUpdateNodeInternals,
} from "@xyflow/react";
import { CircleX, InfoIcon, Pencil } from "lucide-react";
import { ReactNode, useCallback, useId, useState } from "react";
import { useForm, SubmitHandler, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { type GraphTrigger, type GraphState, GraphInstance, TriggerData } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { DialogClose, DialogFooter, DialogHeader } from "@/components/ui/dialog";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import TriggerSelect from "@/components/workflow/trigger-select";
import { z } from "zod";

export default function CustomEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  data,
}: EdgeProps<GraphTrigger>) {
  const reactFlow: GraphInstance = useReactFlow();

  const updateNodeInternals = useUpdateNodeInternals();
  const [edgePath, labelX, labelY] = getSmoothStepPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    borderRadius: 10,
  });
  const deleteEdge = useCallback(
    function () {
      const nodes = reactFlow.getNodes();
      let edges = reactFlow.getEdges();
      let deletedEdge: Edge | null = null;
      for (let index in edges) {
        if (edges[index].id == id) {
          deletedEdge = edges[index];
          break;
        }
      }
      edges = edges.filter((e) => e.id !== (deletedEdge ?? { id: "" }).id);
      for (let index in nodes) {
        if (deletedEdge && nodes[index].id == deletedEdge.source) {
          (nodes[index] as GraphState).data.edges = (nodes[index] as GraphState).data.edges.filter(
            (edge) => edge.id != deletedEdge.sourceHandle
          );
          updateNodeInternals(nodes[index].id);
        }
      }
      reactFlow.setEdges(edges);
      reactFlow.setNodes(nodes);
      if (deletedEdge) {
        updateNodeInternals(deletedEdge.source);
      }
    },
    [id, reactFlow, updateNodeInternals]
  );

  function updateEdgeData(value: TriggerData) {
    reactFlow.updateEdgeData(id, {
      trigger: value.trigger,
      label: value.label,
    });
  }

  return (
    <>
      <BaseEdge id={id} path={edgePath} />
      <EdgeLabelRenderer>
        <div
          style={{
            position: "absolute",
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            pointerEvents: "all",
            display: "flex",
            gap: "12px",
          }}
        >
          <Button variant="ghost" className="nodrag nopan rounded-full p-2" onClick={deleteEdge}>
            <CircleX />
          </Button>
          <p className="text-l">{data!.label}</p>
          <EditDialog asChild onValueChange={updateEdgeData}>
            <Button variant="ghost" className="nodrag nopan rounded-full p-2">
              <Pencil size="1rem" />
            </Button>
          </EditDialog>
        </div>
      </EdgeLabelRenderer>
    </>
  );
}

const editSchema = z.object({
  trigger: z.string().min(1, "Trigger is required"),
  label: z.string().trim().min(1, "Label is required"),
});
type EditData = z.infer<typeof editSchema>;

function EditDialog({
  asChild,
  children,
  value,
  onValueChange = () => {},
}: {
  asChild?: boolean;
  children: ReactNode;

  value?: TriggerData;
  onValueChange: (value: TriggerData) => void;
}) {
  const nameId = useId();
  const eventId = useId();
  const formId = useId();

  const [open, setOpen] = useState(true);

  const {
    handleSubmit,
    register,
    formState: { errors },
    control,
  } = useForm<EditData>({
    resolver: zodResolver(editSchema),
    defaultValues: {
      label: value?.label,
      trigger: value?.trigger,
    },
  });

  function onSubmit(data: EditData) {
    setOpen(false);
    onValueChange(data);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit trigger</DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete your account and remove your
            data from our servers.
          </DialogDescription>
        </DialogHeader>
        <form className="flex flex-col gap-1" id={formId} onSubmit={handleSubmit(onSubmit)}>
          <div className="flex items-center gap-2">
            <Label htmlFor={eventId}>Event</Label>
            <Tooltip>
              <TooltipTrigger>
                <InfoIcon size="1rem" />
              </TooltipTrigger>
              <TooltipContent>
                <p className="text-sm">What event will cause the transition to occur</p>
              </TooltipContent>
            </Tooltip>
          </div>

          <Controller
            name="trigger"
            control={control}
            rules={{ required: true }}
            render={({ field }) => {
              const { onChange, ...field2 } = field;
              return <TriggerSelect id={eventId} onValueChange={onChange} {...field2} />;
            }}
          />
          {errors.trigger && (
            <p role="alert" className="text-red-500 text-sm">
              {errors.trigger.message}
            </p>
          )}

          <div className="flex items-center gap-2 mt-5">
            <Label htmlFor={nameId}>Label</Label>
            <Tooltip>
              <TooltipTrigger>
                <InfoIcon size="1rem" />
              </TooltipTrigger>
              <TooltipContent>
                <p className="text-sm">Human-readable label describing the trigger</p>
              </TooltipContent>
            </Tooltip>
          </div>

          <Input id={nameId} {...register("label")} />

          {errors.label && (
            <p role="alert" className="text-red-500 text-sm">
              {errors.label.message}
            </p>
          )}
        </form>
        <DialogFooter>
          <Button form={formId}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
