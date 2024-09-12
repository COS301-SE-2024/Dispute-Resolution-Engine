"use client";
import { useCallback, useEffect, useId, useMemo, useRef, useState } from "react";
import {
  ReactFlow,
  addEdge,
  useNodesState,
  useEdgesState,
  Background,
  Connection,
  Edge,
} from "@xyflow/react";
import CustomEdge from "./CustomEdge";

import "@xyflow/react/dist/style.css";
import { Button } from "@/components/ui/button";
import { FormProvider, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Textarea } from "@/components/ui/textarea";
import CustomNode from "./CustomNode";

const initialNodes = [
  {
    id: "0",
    type: "customNode",
    position: { x: 0, y: 0 },
    data: { label: "Node A", edges: [{ id: "testId" }] },
  },
  {
    id: "1",
    type: "customNode",
    position: { x: 0, y: 100 },
    data: { label: "Node B", edges: [{ id: "testId2" }] },
  },
  {
    id: "2",
    type: "customNode",
    position: { x: 0, y: 200 },
    data: { label: "Node C", edges: [] },
  },
];

const initialEdges : Edge[]= [
  { id: "0->1", type: "custom-edge", source: "0", target: "1" },
  { id: "1->2", type: "custom-edge", source: "1", target: "2" },
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
  let currId = useRef(3);
  let currEdgeId = useRef(1);
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const nodeTypes = useMemo(() => ({ customNode: CustomNode }), []);
  const onConnect = useCallback(
    (connection: Connection) => {
      const edge = { ...connection, type: "custom-edge"};
      setEdges(
        (eds) =>
          addEdge(edge, eds) as { id: string; type: string; source: string; target: string }[],
      );
      // TODO MAKE THE PAGE manage a latest handle that the node uses for it's new handle
      // TODO morph "new" ids into unique ids in page
      setNodes((node) => {
        console.log("setting nodes", edges, node);
        for (var index in node) {
          var currEdges = [];
          if (connection.source == node[index].id) {
            currEdges.push({ id: currEdgeId.current});
            currEdgeId.current = currEdgeId.current + 1
          }
          for (var edgeIndex in edges) {
            if (edges[edgeIndex].source == node[index].id) {
              currEdges.push({ id: edges[index].sourceHandle ? edges[index].sourceHandle : ""} );
            }
          }
          console.log(currEdges);
          node[index].data.edges = currEdges.map(edge => ({ id: (edge.id?.toString() ?? "") }));
        }
        return node;
      });
    },
    [setEdges, setNodes, edges],
  );

  const addNode = useCallback(
    (params: any) => {
      const newNode = {
        id: currId.current.toString(),
        type: "customNode",
        position: { x: 0, y: 200 },
        data: { label: params.label, edges: [{ id: "hi" }] },
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
    <div className="h-full grid grid-rows-[1fr_auto]">
      <ReactFlow
        className="dark:bg-surface-dark-950"
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        edgeTypes={edgeTypes}
        nodeTypes={nodeTypes}
        colorMode="system"
        fitView
      />
      <form
        onSubmit={form.handleSubmit(addNode)}
        className="p-5 bg-surface-light-50 dark:bg-surface-dark-900 border-t border-primary-500/30"
      >
        <Textarea {...form.register("label")}></Textarea>
        <Button type="submit">ADD NODE</Button>
      </form>
    </div>
  );
}

export default Flow;
