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
import { ReactNode, useCallback, useState } from "react";

import { type GraphTrigger, type GraphState, GraphInstance } from "@/lib/types";
import EditForm from "./edit-form";
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

  /** Used to determine when a component the label of a node is being edited */
  const [editing, setEditing] = useState(false);

  function setEdgeLabel(value: string) {
    setEditing(false);
    reactFlow.updateEdgeData(id, {
      trigger: value,
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
          {editing ? (
            <EditForm
              value={data!.trigger}
              onCommit={setEdgeLabel}
              onCancel={() => setEditing(false)}
            />
          ) : (
            <>
              <Button
                variant="ghost"
                className="nodrag nopan rounded-full p-2"
                onClick={deleteEdge}
              >
                <CircleX />
              </Button>
              <p className="text-l">{data!.trigger}</p>

              <EditDialog asChild>
                <Button variant="ghost" className="nodrag nopan rounded-full p-2">
                  <Pencil size="1rem" />
                </Button>
              </EditDialog>
            </>
          )}
        </div>
      </EdgeLabelRenderer>
    </>
  );
}

function EditDialog({ asChild, children }: { asChild?: boolean; children: ReactNode }) {
  return (
    <Dialog>
      <DialogTrigger asChild={asChild}>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit trigger</DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete your account and remove your
            data from our servers.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <Label>Label</Label>
            <Tooltip>
              <TooltipTrigger>
                <InfoIcon size="1rem" />
              </TooltipTrigger>
              <TooltipContent>
                <p className="text-sm">Human-readable label describing the trigger</p>
              </TooltipContent>
            </Tooltip>
          </div>

          <Input />

          <div className="flex items-center gap-2 mt-5">
            <Label>Event</Label>
            <Tooltip>
              <TooltipTrigger>
                <InfoIcon size="1rem" />
              </TooltipTrigger>
              <TooltipContent>
                <p className="text-sm">What event will cause the transition to occur</p>
              </TooltipContent>
            </Tooltip>
          </div>

          <Input />
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button>Confirm</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
