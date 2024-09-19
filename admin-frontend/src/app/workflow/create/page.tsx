"use client";
import { useCallback, useMemo, useRef } from "react";
import {
  ReactFlow,
  useNodesState,
  useEdgesState,
  Connection,
  useReactFlow,
  ReactFlowProvider,
  useUpdateNodeInternals,
  ConnectionState,
} from "@xyflow/react";
import CustomEdge from "./CustomEdge";

import "@xyflow/react/dist/style.css";
import { Button } from "@/components/ui/button";
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

/**
 * Used for assigning IDs to both nodes an edges. This is required because
 * useId cannot be called inside a useCallback function, so a custom
 * implementation is required.
 */
function useCustomId(start: number | undefined) {
  let currId = useRef(start ?? 0);
  return function () {
    const id = currId.current.toString();
    currId.current++;
    return id;
  };
}

// http://localhost:3000/workflow
function Flow() {
  const createId = useCustomId(initialNodes.length);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const reactFlowInstance: GraphInstance = useReactFlow();

  function createEdge(connection: Connection, trigger: string): GraphTrigger {
    const edge = {
      ...connection,
      id: createId(),
      data: { trigger },
      type: "custom-edge",
    } satisfies GraphTrigger;

    const sourceNode = reactFlowInstance.getNode(connection.source)!;
    sourceNode.data.edges.push({
      id: connection.sourceHandle ?? "whyNoHandle",
    });
    updateNodeInternals(connection.source);

    return edge;
  }

  const nodeTypes = useMemo(() => ({ customNode: CustomNode }), []);
  const updateNodeInternals = useUpdateNodeInternals();
  const onConnect = useCallback(
    (connection: Connection) => {
      if (connection.sourceHandle === "new") {
        connection.sourceHandle = createId();
        reactFlowInstance.addEdges([createEdge(connection, "bruh")]);
      }
    },
    [reactFlowInstance, updateNodeInternals],
  );
  const { screenToFlowPosition } = useReactFlow();
  const onConnectEnd = useCallback(
    (event: any, connectionState: Omit<ConnectionState, "inProgress">) => {
      if (!connectionState.isValid) {
        const { clientX, clientY } = "changedTouches" in event ? event.changedTouches[0] : event;
        const newNode: GraphState = {
          id: createId(),
          type: "customNode",
          position: screenToFlowPosition({
            x: clientX,
            y: clientY,
          }),
          data: { label: "New Node", edges: [] },
        };

        const newEdge: GraphTrigger = createEdge(
          {
            source: connectionState.fromNode?.id ?? "",
            target: newNode.id,
            sourceHandle: createId(),
            targetHandle: null,
          },
          "new_trigger",
        );
        reactFlowInstance.addNodes([newNode]);
        reactFlowInstance.addEdges([newEdge]);
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

  return (
    <div className="h-full grid grid-rows-[1fr_auto]">
      <Flow></Flow>
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
