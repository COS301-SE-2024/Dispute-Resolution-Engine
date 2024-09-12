"use client";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { ReactFlow, addEdge, useNodesState, useEdgesState, Background, Connection } from "@xyflow/react";
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
    data: { label: "Node A", edges: [{id: "testId"}] },
  },
  {
    id: "1",
    type: "customNode",
    position: { x: 0, y: 100 },
    data: { label: "Node B", edges: [{id: "testId2"}] },
  },
  {
    id: "2",
    type: "customNode",
    position: { x: 0, y: 200 },
    data: { label: "Node C", edges: [] },
  },
];

const initialEdges = [
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
  let currId = useRef(3)
  let currEdgeId = useRef(0)
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const nodeTypes = useMemo(() => ({ customNode: CustomNode }), []);
  const onConnect = useCallback(
    (connection: Connection) => {
      const edge = { ...connection, type: "custom-edge" };
      setEdges((eds) => addEdge(edge, eds) as { id: string; type: string; source: string; target: string; }[]);
      setNodes((node) => {
        console.log("setting nodes", edges, node)
        for(var index in node){
          var currEdges = []
          if(connection.source == node[index].id){
            currEdges.push({id: "newConn"})
          }
          for(var edgeIndex in edges){
            if (edges[edgeIndex].source == node[index].id) {
              currEdges.push({id: edges[edgeIndex].target})
            }
          }
          console.log(currEdges)
          node[index].data.edges = currEdges
          // if (node[index].id == connection.source || node[index].id == connection.target){
          //   node[index].data.edges.push({id: "newEdge: " + currEdgeId.current})
          //   currEdgeId.current = currEdgeId.current + 1
          //   console.log("setting correct nodes: ", node[index])
          // }
        }
        return node
      })
    },
    [setEdges, setNodes, edges]
  );

  const addNode = useCallback(
    (params: any) => {
      const newNode = {
        id: currId.current.toString(),
        type: "customNode",
        position: { x: 0, y: 200 },
        data: { label: params.label , edges: [{ id: "hi" }]},
        // data: { label: params.label , time: {hours: 10, minutes: 20, seconds: 30}},
      };
      currId.current = currId.current + 1;
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
        <Background bgColor="#111827" />
      </ReactFlow>
      <form onSubmit={form.handleSubmit(addNode)}>
        <Textarea {...form.register("label")}></Textarea>
        <Button type="submit">ADD NODE</Button>
      </form>
    </div>
  );
}

export default Flow;
