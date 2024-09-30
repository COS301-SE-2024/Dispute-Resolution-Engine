import { Badge } from "@/components/ui/badge";
import { getArchives } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import Link from "next/link";

function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <tr>
      <td className="border px-4 py-2">
        <Link href={`/archive/${props.id}`}>
          <h3 className="hover:underline font-semibold text-lg">{props.title}</h3>
        </Link>
        <p className="dark:text-white/50">{props.summary}</p>
        <div className="space-x-1 mt-2">
          {props.category.map((cat) => (
            <Badge key={cat}>{cat}</Badge>
          ))}
        </div>
      </td>
    </tr>
  );
}

export default async function ArchiveBrowse() {
  let response = await getArchives();
  let data : ArchivedDisputeSummary[];
  if (!response.data) {
    data = [];
  } else {
    data = response.data.archives;
  }

  return (
    <div className="pt-8 pl-8">
      <main className="mx-20">
      <table className="min-w-full rounded-lg overflow-hidden">
        <thead>
        <tr>
          <th className="py-2 text-4xl">Archived Disputes</th>
        </tr>
        </thead>
        <tbody className="bg-dre-500">
        {data.length > 0 ? (
          data.map((dispute) => <SearchResult key={dispute.id} {...dispute} />)
        ) : (
          <tr>
          <td className="border px-4 py-2">No results</td>
          </tr>
        )}
        </tbody>
      </table>
      </main>
    </div>
  );
}