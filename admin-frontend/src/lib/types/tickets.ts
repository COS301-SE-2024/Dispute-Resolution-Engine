import { Filter, Sort } from ".";

export const TICKET_STATUS = ["Open", "Closed", "Solved", "On Hold"] as const;
export type TicketStatus = (typeof TICKET_STATUS)[number];

export interface TicketSummary {
  id: string;

  // The user that created the ticket
  user: {
    id: string;
    full_name: string;
  };

  // Timestamp of when the ticket was created
  date_created: string;

  // Timestamp of when the ticket was created
  subject: string;

  // The status of the ticket
  status: TicketStatus;
}

export interface Ticket extends TicketSummary {
  dispute_id: number;

  // The initial message submitted with the ticket
  body: string;

  // All messages exchanged in the ticket
  messages: TicketMessage[];
}

export interface TicketMessage {
  id: string;

  // The user that submitted the message
  user: {
    id: string;
    full_name: string;
  };

  // The timestamp when the user submitted the ticket
  date_sent: string;

  // The message in the ticket
  message: string;
}

// ---------------------------------------------------------------------------- REQUEST BODIES
export type TicketSort = Sort<
  | "date_created" // The date is was created
  | "subject" // The subject of the ticket
  | "status" // The status of the ticket
  | "user" // The full name of the user;
>;
export type TicketFilter = Filter<"status">;

export interface TicketListRequest {
  // Search term for the title of disputes
  search?: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  sort?: TicketSort;

  // The filters to apply to data
  filter?: TicketFilter[];
}

export interface TicketListResponse {
  tickets: TicketSummary[];
  total: number;
}

export type TicketDetailsResponse = Ticket;

export interface TicketPatchRequest {
  // Changes the status of the ticket to the passed-in value
  status?: TicketStatus;
}

export interface TicketMessageRequest {
  message: string;
}
export type TicketMessageResponse = TicketMessage;
