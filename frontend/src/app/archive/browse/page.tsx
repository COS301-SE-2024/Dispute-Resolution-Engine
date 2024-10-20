import { Badge } from "@/components/ui/badge";
import { getArchives } from "@/lib/api/archive";
import { ArchivedDisputeSummary } from "@/lib/interfaces/archive";
import Link from "next/link";

function SearchResult(props: ArchivedDisputeSummary) {
  return (
    <tr className="rounded-lg overflow-hidden">
      <td className=" border-gray-500 px-4 py-2 text-left">
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

const mockDisputes: ArchivedDisputeSummary[] = [
  {
    id: "1",
    title: "Dispute over Contract Terms",
    summary: "A dispute regarding the terms of a contract between two parties.",
    description: "Detailed description of the dispute over contract terms.",
    category: ["Contract", "Legal"],
    date_filed: "2022-01-01",
    date_resolved: "2022-02-01",
    resolution: "Resolved through mediation.",
  },
  {
    id: "2",
    title: "Property Boundary Dispute",
    summary: "A dispute over the boundary lines of a property.",
    description: "Detailed description of the property boundary dispute.",
    category: ["Property", "Legal"],
    date_filed: "2022-03-01",
    date_resolved: "2022-04-01",
    resolution: "Resolved through arbitration.",
  },
  {
    id: "3",
    title: "Employment Contract Dispute",
    summary: "A dispute regarding the terms of an employment contract.",
    description: "Detailed description of the employment contract dispute.",
    category: ["Employment", "Legal"],
    date_filed: "2022-05-01",
    date_resolved: "2022-06-01",
    resolution: "Resolved through negotiation.",
  },
  {
    id: "4",
    title: "Intellectual Property Dispute",
    summary: "A dispute over the ownership of intellectual property.",
    description: "Detailed description of the intellectual property dispute.",
    category: ["Intellectual Property", "Legal"],
    date_filed: "2022-07-01",
    date_resolved: "2022-08-01",
    resolution: "Resolved through court ruling.",
  },
  {
    id: "5",
    title: "Consumer Rights Dispute",
    summary: "A dispute regarding consumer rights and product quality.",
    description: "Detailed description of the consumer rights dispute.",
    category: ["Consumer Rights", "Legal"],
    date_filed: "2022-09-01",
    date_resolved: "2022-10-01",
    resolution: "Resolved through settlement.",
  },
  {
    id: "6",
    title: "Tenant-Landlord Dispute",
    summary: "A dispute between a tenant and landlord over lease terms.",
    description: "Detailed description of the tenant-landlord dispute.",
    category: ["Real Estate", "Legal"],
    date_filed: "2022-11-01",
    date_resolved: "2022-12-01",
    resolution: "Resolved through mediation.",
  },
  {
    id: "7",
    title: "Business Partnership Dispute",
    summary: "A dispute between business partners over profit sharing.",
    description: "Detailed description of the business partnership dispute.",
    category: ["Business", "Legal"],
    date_filed: "2023-01-01",
    date_resolved: "2023-02-01",
    resolution: "Resolved through arbitration.",
  },
  {
    id: "8",
    title: "Insurance Claim Dispute",
    summary: "A dispute over the settlement of an insurance claim.",
    description: "Detailed description of the insurance claim dispute.",
    category: ["Insurance", "Legal"],
    date_filed: "2023-03-01",
    date_resolved: "2023-04-01",
    resolution: "Resolved through negotiation.",
  },
  {
    id: "9",
    title: "Construction Contract Dispute",
    summary: "A dispute over the terms of a construction contract.",
    description: "Detailed description of the construction contract dispute.",
    category: ["Construction", "Legal"],
    date_filed: "2023-05-01",
    date_resolved: "2023-06-01",
    resolution: "Resolved through court ruling.",
  },
  {
    id: "10",
    title: "Healthcare Service Dispute",
    summary: "A dispute regarding the quality of healthcare services provided.",
    description: "Detailed description of the healthcare service dispute.",
    category: ["Healthcare", "Legal"],
    date_filed: "2023-07-01",
    date_resolved: "2023-08-01",
    resolution: "Resolved through settlement.",
  },
];

export default async function ArchiveBrowse() {
  // let response = await getArchives();
  // let data: ArchivedDisputeSummary[];
  // if (!response.data) {
  //   data = [];
  // } else {
  //   data = response.data.archives;
  // }
  const data = mockDisputes;

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
