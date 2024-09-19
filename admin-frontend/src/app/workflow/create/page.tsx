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
  ConnectionState,
  type ReactFlowInstance,
} from "@xyflow/react";
import CustomEdge from "./CustomEdge";

import "@xyflow/react/dist/style.css";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Textarea } from "@/components/ui/textarea";
import CustomNode from "./CustomNode";

import { type Workflow, type GraphState, type GraphTrigger, GraphInstance } from "@/lib/types";

const initialNodes: GraphState[] = [
  {
    id: "0",
    type: "customNode",
    position: { x: 0, y: 0 },
    data: { label: "Node A", edges: [] },
  },
];

const initialEdges: GraphTrigger[] = [];

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
  const reactFlowInstance = useReactFlow();
  const updateNodeInternals = useUpdateNodeInternals();
  const onConnect = useCallback(
    (connection: Connection) => {
      const edges = reactFlowInstance.getEdges();
      const nodes = reactFlowInstance.getNodes();

      let newEdge = false;
      if (connection.sourceHandle === "new") {
        connection.sourceHandle = currId.current.toString();
        currId.current = currId.current + 1;
        edges.push({
          ...connection,
          data: {
            trigger: "bruh",
          },
          type: "custom-edge",
          id: currId.current.toString(),
        } satisfies GraphTrigger);
        currId.current++;
        newEdge = true;
      }

      for (let nodeIndex in nodes) {
        if (nodes[nodeIndex].id === connection.source) {
          (nodes[nodeIndex] as GraphState).data.edges.push({
            id: connection.sourceHandle ?? "whyNoHandle",
          });
          updateNodeInternals(nodes[nodeIndex].id);
        }
      }
      console.log("connect");
      reactFlowInstance.setEdges(edges);
      reactFlowInstance.setNodes(nodes);
    },
    [reactFlowInstance, updateNodeInternals],
  );
  const { screenToFlowPosition } = useReactFlow();
  const onConnectEnd = useCallback(
    (event: any, connectionState: Omit<ConnectionState, "inProgress">) => {
      if (!connectionState.isValid) {
        const edges = reactFlowInstance.getEdges();
        const nodes = reactFlowInstance.getNodes();
        const nodeId = currId.current.toString();
        currId.current++;
        const handleId = currId.current;
        currId.current++;
        const { clientX, clientY } = "changedTouches" in event ? event.changedTouches[0] : event;
        const newNode: GraphState = {
          id: nodeId,
          type: "customNode",
          position: screenToFlowPosition({
            x: clientX,
            y: clientY,
          }),
          data: { label: "New Node", edges: [] },
        };
        const edgeId = currId.current;
        currId.current++;
        const newEdge: Edge = {
          id: edgeId.toString(),
          source: connectionState.fromNode?.id ?? "",
          target: nodeId,
          data: {
            trigger: "new trigger",
          },
          type: "custom-edge",
          sourceHandle: handleId.toString(),
        };
        for (let index in nodes) {
          if (nodes[index].id == newEdge.source) {
            (nodes[index] as GraphState).data.edges.push({ id: handleId.toString() });
          }
        }
        nodes.push(newNode);
        edges.push(newEdge);
        reactFlowInstance.setNodes(nodes);
        reactFlowInstance.setEdges(edges);
        updateNodeInternals(newNode.id);
      }
    },
    [reactFlowInstance, screenToFlowPosition, updateNodeInternals],
  );

  return (
    <ReactFlow
      className="dark:bg-surface-dark-950 stroke-primary-500"
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      onConnectEnd={onConnectEnd}
      edgeTypes={edgeTypes}
      nodeTypes={nodeTypes}
      colorMode="system"
      fitView
    />
  );
}

function InnerProvider() {
  const currNodeId = useRef(500);
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const reactFlow: GraphInstance = useReactFlow();

  function convertWorkflow() {
    const { nodes, edges } = reactFlow.toObject();
    const workflow: Workflow = {
      label: "bruh",
      initial: "Im not sure",
      states: Object.fromEntries(
        nodes.map((node) => [
          node.id,
          {
            label: node.data.label,
            description: "sure bud",
            events: Object.fromEntries(
              edges
                .filter((edge) => edge.source == node.id)
                .map((edge) => [
                  edge.id,
                  {
                    label: "oi blud, do somfin",
                    next_state: edge.target,
                  },
                ]),
            ),
          },
        ]),
      ),
    };
    console.log(workflow);
  }
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
      <Flow></Flow>
      <form
        onSubmit={form.handleSubmit(addNode)}
        className="p-5 bg-surface-light-50 dark:bg-surface-dark-900 border-t border-primary-500/30"
      >
        <Textarea {...form.register("label")}></Textarea>
        <Button type="submit">ADD NODE</Button>
      </form>
      <Button onClick={convertWorkflow}>Convert to workflow</Button>
    </div>
  );
}

function ProviderBS() {
  return (
    <ReactFlowProvider>
      <InnerProvider />
    </ReactFlowProvider>
  );
}
export default ProviderBS;
