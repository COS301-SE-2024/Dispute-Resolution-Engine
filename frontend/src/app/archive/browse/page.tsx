import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import { SearchResult } from "../search/page";

// export interface ArchivedDisputeSummary {
//     id: string;

//     title: string;
//     summary: string;
//     description: string;

//     category: string[];

//     date_filed: string;
//     date_resolved: string;

//     resolution: string;
//   }
export default function ArchiveBrowse() {
  const data: { archives: ArchivedDisputeSummary[] } = {
    archives: [
      {
        id: "1",
        title: "TITLe",
        summary: "Summar",
        description: "Descfertarsd",
        category: ["asdfsad"],
        date_filed: "satar",
        date_resolved: "endsfdas",
        resolution: "dones",
      },
      {
        id: "2",
        title: "TITLe",
        summary: "Summar",
        description: "Descfertarsd",
        category: ["asdfsad"],
        date_filed: "satar",
        date_resolved: "endsfdas",
        resolution: "dones",
      },
      {
        id: "3",
        title: "TITLe",
        summary: "Summar",
        description: "Descfertarsd",
        category: ["asdfsad"],
        date_filed: "satar",
        date_resolved: "endsfdas",
        resolution: "dones",
      },
    ],
  };
  return (
    <div className="pt-8 pl-8">
      <main className="mx-20 grid grid-cols-2">
        <ol className="space-y-5">
          {data!.archives.length > 0 ? (
            data!.archives.map((dispute) => <SearchResult key={dispute.id} {...dispute} />)
          ) : (
            <p>No results</p>
          )}
        </ol>
      </main>
    </div>
  );
}
