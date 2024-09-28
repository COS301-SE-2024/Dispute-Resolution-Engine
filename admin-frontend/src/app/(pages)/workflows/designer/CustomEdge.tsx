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
import { CircleX, Pencil } from "lucide-react";
import { useCallback, useState } from "react";

import { type GraphTrigger, type GraphState, GraphInstance } from "@/lib/types";
import EditForm from "./edit-form";
import { Button } from "@/components/ui/button";

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
            (edge) => edge.id != deletedEdge.sourceHandle,
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
    [id, reactFlow, updateNodeInternals],
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
              <Button
                variant="ghost"
                className="nodrag nopan rounded-full p-2"
                onClick={() => setEditing(true)}
              >
                <Pencil size="1rem" />
              </Button>
            </>
          )}
        </div>
      </EdgeLabelRenderer>
    </>
  );
}
