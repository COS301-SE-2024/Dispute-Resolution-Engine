import { Badge } from "@/components/ui/badge";
import { getArchives } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import Link from "next/link";
function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <li>
      <div className="flex items-center gap-5 mb-3">
        <Link href={`/archive/${props.id}`}>
          <h3 className="hover:underline font-semibold text-lg">{props.title}</h3>
        </Link>
        <div className="space-x-1">
          {props.category.map((cat) => (
            <Badge key={cat}>{cat}</Badge>
          ))}
        </div>
      </div>
      <p className="dark:text-white/50">{props.summary}</p>
    </li>
  );
}

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
export default async function ArchiveBrowse() {
  // const data: { archives: ArchivedDisputeSummary[] } = {
  //   archives: [
  //     {
  //       id: "1",
  //       title: "TITLe",
  //       summary: "Summar",
  //       description: "Descfertarsd",
  //       category: ["asdfsad"],
  //       date_filed: "satar",
  //       date_resolved: "endsfdas",
  //       resolution: "dones",
  //     },
  //     {
  //       id: "2",
  //       title: "TITLe",
  //       summary: "Summar",
  //       description: "Descfertarsd",
  //       category: ["asdfsad"],
  //       date_filed: "satar",
  //       date_resolved: "endsfdas",
  //       resolution: "dones",
  //     },
  //     {
  //       id: "3",
  //       title: "TITLe",
  //       summary: "Summar",
  //       description: "Descfertarsd",
  //       category: ["asdfsad"],
  //       date_filed: "satar",
  //       date_resolved: "endsfdas",
  //       resolution: "dones",
  //     },
  //   ],
  // };
  const {data, error} = await getArchives()
  console.log(data)
  return (
    <div className="pt-8 pl-8">
      <main className="mx-20 grid grid-cols-2">
        <ol className="space-y-5">
          {data!.length > 0 ? (
            data!.map((dispute) => <SearchResult key={dispute.id} {...dispute} />)
          ) : (
            <p>No results</p>
          )}
        </ol>
      </main>
    </div>
  );
}
