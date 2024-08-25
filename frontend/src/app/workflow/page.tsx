"use client"
import { useCallback } from 'react';
import {
  ReactFlow,
  addEdge,
  useNodesState,
  useEdgesState,
} from '@xyflow/react';
import CustomEdge from './CustomEdge';

import '@xyflow/react/dist/style.css';
import { Button } from '@/components/ui/button';

const initialNodes = [
  { id: 'a', position: { x: 0, y: 0 }, data: { label: 'Node A' } },
  { id: 'b', position: { x: 0, y: 100 }, data: { label: 'Node B' } },
  { id: 'c', position: { x: 0, y: 200 }, data: { label: 'Node C' } },
];

const initialEdges = [
  { id: 'a->b', type: 'custom-edge', source: 'a', target: 'b' },
  { id: 'b->c', type: 'custom-edge', source: 'b', target: 'c' },
];

const edgeTypes = {
  'custom-edge': CustomEdge,
};

function Flow() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const onConnect = useCallback(
    (connection : any) => {
      const edge = { ...connection, type: 'custom-edge' };
      setEdges((eds) => addEdge(edge, eds));
    },
    [setEdges],
  );
  const addNode  = useCallback(() => {
    const newNode  = { id: 'd', position: { x: 0, y: 200 }, data: { label: 'NEW NODE YIPPEE' } }
    setNodes((nds) => nds.concat(newNode));  
  
  }, []);
  return (
    <div className="h-96">
    <ReactFlow className="h-24"
      nodes={nodes}
      edges={edges}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      edgeTypes={edgeTypes}
      colorMode='dark'
      fitView
    >

    </ReactFlow>
    <Button onClick={addNode}>ADD NODE</Button>
    </div>
  );
}

export default Flow;