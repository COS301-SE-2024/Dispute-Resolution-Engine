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
import { CircleX } from "lucide-react";
import { useCallback } from "react";

import { type GraphTrigger, type GraphState } from "@/lib/types";

export default function CustomEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  data,
}: EdgeProps<GraphTrigger>) {
  const { setEdges, getEdges } = useReactFlow();
  const { setNodes } = useReactFlow();
  const reactFlowInstance = useReactFlow();
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
      const nodes = reactFlowInstance.getNodes();
      let edges = reactFlowInstance.getEdges();
      console.log("edges before ", edges);
      console.log("nodes before ", nodes);
      console.log("id before ", id);
      let deletedEdge: Edge | null = null;
      for (let index in edges) {
        if (edges[index].id == id) {
          deletedEdge = edges[index];
          break;
        }
      }
      console.log("deletedEdge ", deletedEdge);
      edges = edges.filter((e) => e.id !== (deletedEdge ?? { id: "" }).id);
      for (let index in nodes) {
        if (deletedEdge && nodes[index].id == deletedEdge.source) {
          (nodes[index] as GraphState).data.edges = (nodes[index] as GraphState).data.edges.filter(
            (edge) => edge.id != deletedEdge.sourceHandle,
          );
          updateNodeInternals(nodes[index].id);
        }
      }
      console.log("edges after ", edges);
      console.log("nodes after ", nodes);
      reactFlowInstance.setEdges(edges);
      reactFlowInstance.setNodes(nodes);
      if (deletedEdge) {
        updateNodeInternals(deletedEdge.source);
      }
    },
    [id, reactFlowInstance, updateNodeInternals],
  );
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
          <h1 contentEditable="true" className="text-l" suppressContentEditableWarning={true}>
            {data?.trigger ?? "error"}
          </h1>
          <button className="nodrag nopan" onClick={deleteEdge}>
            <CircleX />
          </button>
        </div>
      </EdgeLabelRenderer>
    </>
  );
}
