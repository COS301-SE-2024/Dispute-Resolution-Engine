"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
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

import {
  type GraphState,
  type GraphTrigger,
  type GraphInstance,
  WorkflowCreateRequest,
  Workflow,
} from "@/lib/types";
import {
  createWorkflow,
  graphToWorkflow,
  updateWorkflow,
  workflowToGraph,
} from "@/lib/api/workflow";
import { Textarea } from "@/components/ui/textarea";
import WorkflowTitle from "@/components/workflow/workflow-title";
import { SaveIcon } from "lucide-react";
import { useCustomId } from "@/lib/hooks/use-customid";
import { useToast } from "@/lib/hooks/use-toast";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { usePathname, useRouter } from "next/navigation";

const initialNodes: GraphState[] = [
  {
    id: "0",
    type: "customNode",
    position: { x: 0, y: 0 },
    data: { initial: true, label: "New state", description: "Initial state", edges: [] },
  },
];

const initialEdges: GraphTrigger[] = [];

const edgeTypes = {
  "custom-edge": CustomEdge,
};

// http://localhost:3000/workflow
function Flow({ setIsSaved }: { setIsSaved: any }) {
  const createId = useCustomId(initialNodes.length);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);

  const handleNodesChange = useCallback(
    (changes: any) => {
      onNodesChange(changes);
      setIsSaved(false);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [onNodesChange]
  );
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const handleEdgesChange = useCallback(
    (changes: any) => {
      onNodesChange(changes);
      setIsSaved(false);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [onEdgesChange]
  );
  const reactFlowInstance: GraphInstance = useReactFlow();
  function createEdge(connection: Connection, trigger: string): GraphTrigger {
    const edge = {
      ...connection,
      id: createId(),
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
        updateNodeInternals(connection.source);
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
          data: { label: "New state", edges: [], description: "New state" },
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
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [reactFlowInstance, screenToFlowPosition, updateNodeInternals]
  );

  return (
    <ReactFlow
      className="dark:bg-surface-dark-950 stroke-primary-500"
      nodes={nodes}
      edges={edges}
      onNodesChange={handleNodesChange}
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

function InnerPage({ workflow }: { workflow?: Workflow }) {
  useEffect(() => {
    if (!workflow) {
      return;
    }
    setWorkflow(workflow);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [workflow]);

  const reactFlow: GraphInstance = useReactFlow();
  async function setWorkflow(wf: Workflow) {
    const [nodes, edges] = await workflowToGraph(wf.definition);
    let idTrack: number = 100;
    const dagreGraph = new dagre.graphlib.Graph();
    const nodeWidth = 200;
    const nodeHeight = 100;
    dagreGraph.setDefaultEdgeLabel(() => ({}));
    dagreGraph.setGraph({ rankdir: "LR", ranker: "tight-tree" });
    nodes.forEach((node) => {
      dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
    });
    edges.forEach((edge) => {
      dagreGraph.setEdge(edge.source, edge.target);
    });
    dagre.layout(dagreGraph);
    for (let node of nodes) {
      const nodeWithPosition = dagreGraph.node(node.id);
      node.position = {
        x: nodeWithPosition.x,
        y: nodeWithPosition.y,
      };
    }

    for (let edge of edges) {
      let sourceNode = nodes.find((node) => node.id === edge.source);
      let currHandleId: string = (idTrack++).toString();
      sourceNode?.data.edges.push({ id: currHandleId });
      edge.sourceHandle = currHandleId;
    }
    reactFlow.setNodes(nodes);
    reactFlow.setEdges(edges);
  }

  const { toast } = useToast();
  const updateNodeInternals = useUpdateNodeInternals();
  const [result, setResult] = useState("");
  const [error, setError] = useState<string>();

  const [isSaved, setIsSaved] = useState<boolean>(false);

  async function toWorkflow() {
    const workflow = await graphToWorkflow(reactFlow.toObject());
    setResult(JSON.stringify(workflow, null, 2));
    setError(undefined);
  }

  const router = useRouter();
  const pathname = usePathname();
  async function saveWorkflow() {
    const definition = await graphToWorkflow(reactFlow.toObject());
    // definition.initial = Object.keys(definition.states)[0];
    console.log(definition);

    console.log(definition);
    try {
      if (workflow) {
        await updateWorkflow(workflow.id, {
          name: title,
          definition,
        });
      } else {
        const workflow = await createWorkflow({ name: title, definition });
        router.replace(`${pathname}?id=${workflow.id}`);
      }
      setIsSaved(true);
    } catch (e: unknown) {
      const error = e as Error;
      toast({
        variant: "error",
        title: "Failed saving workflow",
        description: error?.message,
      });
    }
  }

  const [title, setTitle] = useState(workflow?.name ?? "New Workflow");
  async function commitTitle(value: string) {
    setTitle(value);
    if (workflow) {
      await updateWorkflow(workflow.id, {
        name: value,
      });
    }
  }

  return (
    <div className="h-full grid grid-cols-[1fr_3fr] grid-rows-[auto_1fr]">
      <div className="col-span-2 border-b dark:border-primary-500/30 border-primary-500/20 flex items-center gap-2">
        <WorkflowTitle value={title} onValueChange={commitTitle} />
        <Button variant="ghost" title="Save" onClick={saveWorkflow}>
          <SaveIcon size="1.2rem" />
        </Button>
        <span className={isSaved ? "opacity-100 text-sm" : "opacity-50 text-sm"}>
          {isSaved ? "Saved" : "Unsaved"}
        </span>
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
        </div>
      </div>
      <Flow setIsSaved={setIsSaved}></Flow>
    </div>
  );
}

export default function ClientPage({ workflow }: { workflow?: Workflow }) {
  const [client] = useState(new QueryClient());
  return (
    <QueryClientProvider client={client}>
      <ReactFlowProvider>
        <InnerPage workflow={workflow} />
      </ReactFlowProvider>
    </QueryClientProvider>
  );
}
