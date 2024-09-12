import { Input } from "@/components/ui/input";
import {
  BaseEdge,
  EdgeLabelRenderer,
  getBezierPath,
  getStraightPath,
  useReactFlow,
} from "@xyflow/react";
import { CircleX } from "lucide-react";

export default function CustomEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
}: {
  id: string;
  sourceX: number;
  sourceY: number;
  targetX: number;
  targetY: number;
}) {
  const { setEdges, getEdges } = useReactFlow();
  const { setNodes } = useReactFlow();
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
  });

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
            gap: "12px"
          }}
        >
          <h1 contentEditable="true" className="text-3xl" suppressContentEditableWarning={true}>
            trigger
          </h1>
          <button
            className="nodrag nopan"
            onClick={() => {
              setEdges((es) => es.filter((e) => e.id !== id));
              setNodes((node) => {
                let edges = getEdges();
                edges = edges.filter((e) => e.id !== id);
                console.log("setting nodes", edges, node);
                for (var index in node) {
                  var currEdges = [];
                  for (var edgeIndex in edges) {
                    if (edges[edgeIndex].source == node[index].id) {
                      currEdges.push({ id: edges[edgeIndex].target });
                    }
                  }
                  console.log(currEdges);
                  node[index].data.edges = currEdges;
                }
                return node;
              });
            }}
          >
            <CircleX />
          </button>
        </div>
      </EdgeLabelRenderer>
    </>
  );
}
