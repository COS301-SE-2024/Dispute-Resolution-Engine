"use client";
import { useCallback, useMemo, useRef, useState } from "react";
import { ReactFlow, addEdge, useNodesState, useEdgesState, Background } from "@xyflow/react";
import CustomEdge from "./CustomEdge";

import "@xyflow/react/dist/style.css";
import { Button } from "@/components/ui/button";
import { FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Textarea } from "@/components/ui/textarea";
import CustomNode from "./CustomNode";

const initialNodes = [
  { id: "a", type:"customNode", position: { x: 0, y: 0 }, data: { label: "Node A" } },
  { id: "b", position: { x: 0, y: 100 }, data: { label: "Node B" } },
  { id: "c", position: { x: 0, y: 200 }, data: { label: "Node C" } },
];

const initialEdges = [
  { id: "a->b", type: "custom-edge", source: "a", target: "b" },
  { id: "b->c", type: "custom-edge", source: "b", target: "c" },
];

const edgeTypes = {
  "custom-edge": CustomEdge,
};


const newNodeSchema = z.object({
  label: z.string().min(1).max(50),
});
type NewNodeData = z.infer<typeof newNodeSchema>;

// http://localhost:3000/workflow
function Flow() {
  let currId = useRef(1);
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const nodeTypes = useMemo(() => ({ customNode: CustomNode }), []);
  const onConnect = useCallback(
    (connection: any) => {
      const edge = { ...connection, type: "custom-edge" };
      setEdges((eds) => addEdge(edge, eds));
    },
    [setEdges],
  );

  const addNode = useCallback(
    (params: any) => {
      const newNode = {
        id: currId.current.toString(),
        position: { x: 0, y: 200 },
        data: { label: params.label },
      };
      currId.current = currId.current + 1;
      setNodes((nds) => nds.concat(newNode));
    },
    [setNodes],
  );

  const form = useForm<NewNodeData>({
    defaultValues: {
      label: "New Node",
    },
    resolver: zodResolver(newNodeSchema),
  });

  return (
    <div className="h-96">
      <ReactFlow
        className="h-24 "
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        edgeTypes={edgeTypes}
        nodeTypes={nodeTypes}
        colorMode="dark"
        fitView
      >
        <Background  bgColor="#111827"/>
      </ReactFlow>
      <form onSubmit={form.handleSubmit(addNode)}>
        <Textarea {...form.register("label")}></Textarea>
        <Button type="submit">ADD NODE</Button>
      </form>
    </div>
  );
}

export default Flow;
