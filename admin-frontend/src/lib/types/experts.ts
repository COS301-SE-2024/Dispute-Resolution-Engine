export const OBJECTION_STATUS = ["Review", "Overruled", "Sustained"] as const;
export type ObjectionStatus = (typeof OBJECTION_STATUS)[number];

export interface Objection {
  /** ID of the objection itself. */
  id: number;

  /** ID  of the ticket the objection is related to */
  ticket_id: number;

  /** The full name of the expert being objected to */
  expert_name: string;

  /** The full name of the user that submitted the objection */
  user_name: string;

  /** When the objection was submitted */
  date_submitted: string;

  /** The status of the objection */
  status: ObjectionStatus;
}

// ----------------------------------------------------------------- REQUEST/RESPONSE TYPES
export type ObjectionListResponse = Objection[];

export interface ObjectionListRequest {
  status?: ObjectionStatus;
}
