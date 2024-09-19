import { type Node, type Edge, type ReactFlowInstance } from "@xyflow/react";

export type disputeDuration = {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
};

export type eventType = {
  id: string;
};

export type GraphState = Node<
  {
    edges: { id: string }[];
    label?: any;
  },
  "customNode"
>;

export type GraphTrigger = Edge<
  {
    trigger: string;
  },
  "custom-edge"
>;

export type GraphInstance = ReactFlowInstance<GraphState, GraphTrigger>;

export interface Workflow {
  label: string;
  initial: string;
  states: {
    [key: string]: State;
  };
}

export interface State {
  label: string;
  description: string;
  events: {
    [key: string]: Event;
  };
  timer?: Timer;
}

export interface Timer {
  duration: string;
  on_expire: string;
}

export interface Event {
  label: string;
  next_state: string;
}
