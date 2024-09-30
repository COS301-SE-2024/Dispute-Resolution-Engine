import { Badge } from "@/components/ui/badge";
import { getArchives } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import Link from "next/link";

function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <tr className="border rounded-lg overflow-hidden">
      <td className="border px-4 py-2 text-center">
        <Link href={`/archive/${props.id}`}>
          <h3 className="hover:underline font-semibold text-lg">{props.title}</h3>
        </Link>
        <p>{props.summary}</p>
        <div className="space-x-1 mt-2">
          {props.category.map((cat) => (
            <Badge key={cat}>{cat}</Badge>
          ))}
        </div>
      </td>
      <td className="border px-4 py-2 text-center">
        <p><strong>Date Filed:</strong> {props.date_filed}</p>
        <p><strong>Date Resolved:</strong> {props.date_resolved ?? "-"}</p>
        <p><strong>Resolution:</strong> {props.resolution}</p>
      </td>
    </tr>
  );
}

export default async function ArchiveBrowse() {
  let response = await getArchives();
  let data: ArchivedDisputeSummary[];
  if (!response.data) {
    data = [];
  } else {
    data = response.data.archives;
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4 text-center">Archived Disputes</h1>
      <main className="mx-20">
        <table className="min-w-full divide-y rounded-lg overflow-hidden my-4 bg-dre-300">
          <thead>
            <tr>
              <th className="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider">
                Title & Summary
              </th>
              <th className="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider">
                Dates & Resolution
              </th>
            </tr>
          </thead>
          <tbody>
            {data.length > 0 ? (
              data.map((dispute) => <SearchResult key={dispute.id} {...dispute} />)
            ) : (
              <tr>
                <td className="border px-4 py-2 text-center" colSpan={2}>No results</td>
              </tr>
            )}
          </tbody>
        </table>
      </main>
    </div>
  );
}