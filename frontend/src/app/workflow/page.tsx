"use client"
import { Background, BaseEdge, Controls, getStraightPath, ReactFlow } from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import RemovableEdge from "./RemovableEdge";
const initialNodes = [
  { id: "1", position: { x: 0, y: 0 }, data: { label: "Dispute Created" } },
  { id: "2", position: { x: 0, y: 100 }, data: { label: "Response Due" } },
  { id: "3", position: { x: 0, y: 200 }, data: { label: "Default Judgement" } },
];
const initialEdges = [
  { id: "e1-2", source: "1", target: "2" },
  { id: "e2-3", source: "2", target: "3", type: "removable-edge" },
];
const edgeTypes = {
  "removable-edge": RemovableEdge,
};

export default function Workflow() {
  return (
    <div className="w-full h-full">
      <ReactFlow colorMode="system" nodes={initialNodes} edges={initialEdges} edgeTypes={edgeTypes}>
        <Controls />
      </ReactFlow>
    </div>
  );
}
