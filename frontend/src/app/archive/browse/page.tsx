import { Badge } from "@/components/ui/badge";
import { getArchives } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import Link from "next/link";

function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <tr className="rounded-lg overflow-hidden hover:bg-gray-200 dark:hover:bg-blue-600">
      <td className="border-gray-500 px-4 py-2 text-left">
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
      <td className="border-gray-500 px-8 py-2 text-left">
        <p>
          <strong>Date Filed:</strong> {props.date_filed}
        </p>
        <p>
          <strong>Date Resolved:</strong> {props.date_resolved ?? "-"}
        </p>
        <p>
          <strong>Resolution:</strong> {props.resolution}
        </p>
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
    <div className="p-8 max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold mb-8 text-center text-gray-800 dark:text-gray-100">
        Archived Disputes
      </h1>
      <main className="mx-auto">
        <table className="min-w-full  bg-white dark:bg-blue-950 rounded-xl shadow-lg overflow-hidden">
          <thead className="bg-gray-50 dark:bg-blue-900">
            <tr>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600 dark:text-gray-200 uppercase tracking-wider">
                Title & Summary
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-600 dark:text-gray-200 uppercase tracking-wider">
                Dates & Resolution
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-600 dark:divide-gray-700">
            {data.length > 0 ? (
              data.map((dispute) => <SearchResult key={dispute.id} {...dispute} />)
            ) : (
              <tr>
                <td className="px-6 py-4 text-center text-gray-500 dark:text-gray-400" colSpan={2}>
                  No disputes found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </main>
    </div>
  );
}
