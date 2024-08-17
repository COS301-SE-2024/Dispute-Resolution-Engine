export interface ArchivedDisputeSummary {
  id: string;

  title: string;
  summary: string;
  description: string;
  
  category: string[];

  date_filed: string;
  date_resolved: string;

  resolution: string;
}
export interface ArchivedDispute extends ArchivedDisputeSummary {
  events: {
    timestamp: string;
    type: string;
    description: string;
  }[];
}

// /disputes/archive/search?q=Search&limit=10&sort=asc
export type SortAttribute = "title" | "date_filed" | "date_resolved" | "date_filed" | "time_taken";

export interface ArchiveSearchRequest {
  search: string;

  // Pagination parameters
  limit?: number;
  offset?: number;

  order?: "asc" | "desc";

  // What attribute to sort by
  sort?: SortAttribute;

  filter?: {
    category?: string[];
    time?: number;
  };
}
export type ArchiveSearchResponse = {
  archives: ArchivedDisputeSummary[],
  total: number
};
export type ArchiveGetResponse = ArchivedDispute;

// /disputes/archive/{id} <---- Archive
