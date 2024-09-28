"use client";
import { useCallback, useMemo, useRef, useState } from "react";
import dagre from "dagre";
import {
  ReactFlow,
  useNodesState,
  useEdgesState,
  Connection,
  useReactFlow,
  ReactFlowProvider,
  useUpdateNodeInternals,
  ConnectionState,
  Position,
} from "@xyflow/react";
import CustomEdge from "./CustomEdge";

import "@xyflow/react/dist/style.css";
import { Button } from "@/components/ui/button";
import CustomNode from "./CustomNode";

import { type GraphState, type GraphTrigger, type GraphInstance, WorkflowCreateRequest } from "@/lib/types";
import { createWorkflow, graphToWorkflow, workflowToGraph } from "@/lib/api/workflow";
import { workflowSchema } from "@/lib/schema/workflow";
import { Textarea } from "@/components/ui/textarea";
import WorkflowTitle from "@/components/workflow/workflow-title";
import { SaveIcon } from "lucide-react";

const initialNodes: GraphState[] = [
  {
    id: "0",
    type: "customNode",
    position: { x: 0, y: 0 },
    data: { label: "New Node", edges: [] },
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
    [reactFlowInstance, updateNodeInternals]
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
          "new_trigger"
        );
        reactFlowInstance.addNodes([newNode]);
        reactFlowInstance.addEdges([newEdge]);
      }
    },
    [reactFlowInstance, screenToFlowPosition, updateNodeInternals]
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
  const updateNodeInternals = useUpdateNodeInternals()
  const [result, setResult] = useState("");
  const [error, setError] = useState<string>();
  const [title, setTitle] = useState<string>("New Workflow");

  async function toWorkflow() {
    const workflow = await graphToWorkflow(reactFlow.toObject());
    const tempWorkflow = {...workflow, label: "asdf"}
    setResult(JSON.stringify(tempWorkflow, null, 2));
    setError(undefined);
  }
  async function saveWorkflow() {
    const workflow = await graphToWorkflow(reactFlow.toObject());
    workflow.initial = Object.keys(workflow.states)[0]
    const wfRequest : WorkflowCreateRequest = {
      name: title,
      definition : workflow,
    }
    const response = await createWorkflow(wfRequest)
  }

  async function fromWorkflow() {
    let json: string;
    try {
      json = JSON.parse(result);
    } catch (e) {
      setError((e as Error).message);
      return;
    }
    const { data, error } = workflowSchema.safeParse(json);
    if (error) {
      setError(error.issues[0].message);
      return;
    }

    setError(undefined);
    const [nodes, edges] = await workflowToGraph(data);
    let idTrack : number = 100
    const dagreGraph = new dagre.graphlib.Graph();
    const nodeWidth = 200
    const nodeHeight = 100
    dagreGraph.setDefaultEdgeLabel(() => ({}));
    dagreGraph.setGraph({ rankdir: "LR" , ranker:"tight-tree"});
    nodes.forEach((node) => {
      dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
    });
    edges.forEach((edge) => {
      dagreGraph.setEdge(edge.source, edge.target);
    });
    dagre.layout(dagreGraph);
    for (let node of nodes){
      const nodeWithPosition = dagreGraph.node(node.id);
      node.position = {
        x: nodeWithPosition.x,
        y: nodeWithPosition.y,
      };
    }

    for(let edge of edges){
      let sourceNode = nodes.find(node => node.id === edge.source)
      let currHandleId : string = (idTrack++).toString()
      sourceNode?.data.edges.push({id: currHandleId})
      edge.sourceHandle = currHandleId
    }
    reactFlow.setNodes(nodes);
    reactFlow.setEdges(edges);
  }

  return (
    <div className="h-full grid grid-cols-[1fr_3fr] grid-rows-[auto_1fr]">
      <div className="col-span-2 border-b dark:border-primary-500/30 border-primary-500/20 flex items-center gap-2">
        <WorkflowTitle value={title} onValueChange={(value) => {
          setTitle(value)
        }} />
        <Button variant="ghost" title="Save" onClick={saveWorkflow}>
          <SaveIcon size="1.2rem" />
        </Button>
        <span className="opacity-50 text-sm">Unsaved</span>
      </div>
      <div className="p-2 space-y-2 flex flex-col">
        <Textarea
          className="grow resize-none font-mono"
          value={result}
          onChange={(e) => setResult(e.target.value)}
        />
        <div className="flex flex-col gap-2">
          {error && (
            <p role="alert" className="text-red-500">
              {error}
            </p>
          )}
          <Button onClick={toWorkflow}>Convert graph to workflow</Button>
          <Button onClick={fromWorkflow}>Convert workflow to graph</Button>
        </div>
      </div>
      <Flow></Flow>
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
