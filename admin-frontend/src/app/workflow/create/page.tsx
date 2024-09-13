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
  useReactFlow,
  ReactFlowProvider,
  useUpdateNodeInternals,
} from "@xyflow/react";
import CustomEdge from "./CustomEdge";

import "@xyflow/react/dist/style.css";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Textarea } from "@/components/ui/textarea";
import CustomNode, { CustomNodeType } from "./CustomNode";

const initialNodes = [
  {
    id: "0",
    type: "customNode",
    position: { x: 0, y: 0 },
    data: { label: "Node A", edges: [] },
  },
  {
    id: "1",
    type: "customNode",
    position: { x: 0, y: 100 },
    data: { label: "Node B", edges: [] },
  },
  {
    id: "2",
    type: "customNode",
    position: { x: 0, y: 200 },
    data: { label: "Node C", edges: [] },
  },
  {
    id: "3",
    type: "customNode",
    position: { x: 0, y: 300 },
    data: { label: "Node D", edges: [] },
  },
];

const initialEdges: Edge[] = [
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
  // const onConnect = useCallback(
  //   (connection: Connection) => {
  //     if (connection.sourceHandle == "new") {
  //       connection.sourceHandle = currId.current.toString()
  //       currId.current = currId.current + 1
  //     }
  //     const edge = { ...connection, type: "custom-edge"};
  //     setEdges(
  //       (eds) =>
  //         addEdge(edge, eds) as { id: string; type: string; source: string; target: string }[],
  //     );
  //     // TODO MAKE THE PAGE manage a latest handle that the node uses for it's new handle
  //     // TODO morph "new" ids into unique ids in page
  //     setNodes((node) => {
  //       console.log("setting nodes", edges, node);
  //       for (var index in node) {
  //         var currEdges = [];
  //         if (connection.sourceHandle === currId.current.toString() &&  connection.source === node[index].id) {
  //           currEdges.push({ id: connection.sourceHandle});
  //         }
  //         for (var edgeIndex in edges) {
  //           if (edges[edgeIndex].source == node[index].id) {
  //             currEdges.push({ id: edges[index].sourceHandle ? edges[index].sourceHandle : ""} );
  //           }
  //         }
  //         console.log(currEdges);
  //         node[index].data.edges = currEdges
  //       }
  //       return node;
  //     });
  //   },
  //   [setEdges, setNodes, edges],
  // );
  const reactFlowInstance = useReactFlow()
  const updateNodeInternals = useUpdateNodeInternals()
  const onConnect = useCallback(
    (connection: Connection) => {
      const edges = reactFlowInstance.getEdges();
      const nodes = reactFlowInstance.getNodes();
      console.log("edges before ", edges)
      console.log("nodes before ", nodes)
      let newEdge = false;
      if (connection.sourceHandle === "new") {
        connection.sourceHandle = currId.current.toString();
        currId.current = currId.current + 1;
        edges.push({ ...connection, type: "custom-edge", id: currId.current.toString()} as Edge);
        currId.current++
        newEdge = true;
      }
      for (let nodeIndex in nodes) {
        if (nodes[nodeIndex].id === connection.source) {
          (nodes[nodeIndex] as CustomNodeType).data.edges.push({
            id: connection.sourceHandle ?? "whyNoHandle",
          });
          updateNodeInternals(nodes[nodeIndex].id)
        }
      }
      reactFlowInstance.setEdges(edges);
      reactFlowInstance.setNodes(nodes);
      console.log("edges after ", edges)
      console.log("nodes after ", nodes)
    },
    [reactFlowInstance]
  );

  return (
    <ReactFlow
      className="dark:bg-surface-dark-950 stroke-primary-500"
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
  );
}
function ProviderBS() {
  const currNodeId = useRef(500);
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const addNode = useCallback(
    (params: any) => {
      const newNode = {
        id: currNodeId.current.toString(),
        type: "customNode",
        position: { x: 0, y: 200 },
        data: { label: params.label, edges: [] },
      };
      currNodeId.current = currNodeId.current + 1;
      setNodes((nds) => nds.concat(newNode));
    },
    [setNodes]
  );
  const form = useForm<NewNodeData>({
    defaultValues: {
      label: "New Node",
    },
    resolver: zodResolver(newNodeSchema),
  });
  return (
    <div className="h-full grid grid-rows-[1fr_auto]">
      <ReactFlowProvider>
        <Flow></Flow>
      </ReactFlowProvider>
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
export default ProviderBS;
