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
export interface TicketMessageRequest {
  message: string;
}
export type TicketMessageResponse = TicketMessage;

export interface TicketCreateRequest {
  // The subject of the ticket
  subject: string;

  // The body of the ticket
  body: string;
}

export type TicketCreateResponse = Omit<Ticket, "user">;
