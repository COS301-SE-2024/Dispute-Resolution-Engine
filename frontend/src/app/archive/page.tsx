import Loader from "@/components/Loader";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { fetchArchiveHighlights } from "@/lib/api/archive";
import { ArchivedDisputeSummary, ArchiveSearchResponse } from "@/lib/interfaces/archive";
import { ExternalLink } from "lucide-react";
import { Metadata } from "next";
import Link from "next/link";

export const metadata: Metadata = {
  title: "DRE - Archive",
};

function ArchivedDispute(props: ArchivedDisputeSummary) {
  return (
    <Card className="flex flex-col max-w-sm">
      <CardHeader className="flex flex-row items-center justify-between flex-wrap">
        <div>
          <CardTitle>{props.title}</CardTitle>
          <p>{props.date_resolved}</p>
        </div>
        <ul className="flex gap-1 flex-wrap">
          {props.category.map((cat, i) => (
            <Badge key={i}>{cat}</Badge>
          ))}
        </ul>
      </CardHeader>
      <CardContent asChild className="dark:text-white/50 grow truncate">
        <p>{props.summary}</p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button asChild>
          <Link href={`/archive/${props.id}`}>
            <ExternalLink size="1rem" className="mr-2" />
            Read More
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

const mockArchivedDisputeSummaries: ArchivedDisputeSummary[] = [
  {
    id: "1",
    title: "Refund Request Dispute",
    summary: "Customer requested a refund for a damaged product.",
    description: "A customer filed a dispute requesting a refund after receiving a damaged item in their order. The case involved review of the product condition and purchase history.",
    category: ["Refund", "Product Issue"],
    date_filed: "2023-04-12",
    date_resolved: "2023-05-01",
    resolution: "Full refund was issued to the customer.",
  },
  {
    id: "2",
    title: "Service Agreement Dispute",
    summary: "Dispute over terms of a servicadsfjhkojfdsajhkfdsjhkldfsajhklfdahsjkljkhldfsajkhlasdfjkhldsafkjlhasdfhljkfdasjkhlkjhasdfljkhdsfajklhasfdkljhsadflkjhadsflkjhe agreement.",
    description: "A client disputed the terms ofjahksdfjhasdflhkjfdsajkhdfsajkhlfadsjhkfadsjhklfdsajhlkfadslhjkfadshjkldfaskljhdfslhkjdfsjlkhfdslhkjfdsajkhfdskjhfdaskljhfdsajkhlafdsljkhsfdakjhladsfljkhfsdlhjkfdshljkdfsaljkhadfskjlhsdfalkjhfadslhkjfdsalkjhasfdlkhjasdflkjhafsd a service agreement, claiming that certain services were not provided as outlined in the contract. This led to a review of the contract and communications.",
    category: ["Contract", "Service"],
    date_filed: "2023-02-15",
    date_resolved: "2023-03-03",
    resolution: "Mutual agreement reached; contract revised.",
  },
  {
    id: "3",
    title: "Unauthorized Charge Dispute",
    summary: "Customer disputed an unauthorized charge on their account.",
    description: "The customer reported an unauthorized charge on their account for a service they didn't request. Investigation into transaction logs and communications followed.",
    category: ["Billing", "Unauthorized Transaction"],
    date_filed: "2023-06-20",
    date_resolved: "2023-06-30",
    resolution: "Unauthorized charge was reversed, and additional fraud prevention measures were implemented.",
  },
  {
    id: "4",
    title: "Product Warranty Dispute",
    summary: "Dispute regarding the terms of a product warranty.",
    description: "A customer filed a dispute claiming that their product warranty should cover certain repairs that the company initially refused. The terms of the warranty and the product's condition were reviewed.",
    category: ["Warranty", "Product Issue"],
    date_filed: "2023-01-08",
    date_resolved: "2023-01-25",
    resolution: "Repairs were covered under the warranty after review.",
  },
  {
    id: "5",
    title: "Subscription Cancellation Dispute",
    summary: "Customer filed a dispute about issues with canceling their subscription.",
    description: "The customer experienced issues canceling their subscription through the company's platform, resulting in ongoing charges. The case involved a review of the cancellation process.",
    category: ["Subscription", "Billing"],
    date_filed: "2023-05-15",
    date_resolved: "2023-06-01",
    resolution: "Subscription was successfully canceled, and charges were refunded.",
  }
];

export default async function Archive() {
  // const { data, error } = await fetchArchiveHighlights(3);
  // if (error) {
  //   return <h1>{error}</h1>;
  // }
  const mockData: ArchiveSearchResponse = {
    archives: mockArchivedDisputeSummaries,
    total: 5
  }
  const data = mockData
  return (
    <div className="flex flex-col items-center justify-center h-full gap-5 w-2/3 mx-auto">
      <header className="mx-auto w-fit text-center">
        <h1 className="text-6xl font-bold tracking-wide">Archive</h1>
        <p className="dark:text-white/50">Explore our previously handled cases</p>
      </header>
      <main className="w-2/3">
        <form action="/archive/search" className="flex flex-col items-center gap-2">
          <Input
            name="q"
            className="rounded-full dark:bg-dre-bg-light/5 px-6 py-4 border-none"
            placeholder="Search the Archive..."
          />
          <div className="flex gap-2">
            <Button type="submit">Search</Button>
            <Button asChild>
              <Link href="/archive/browse">
                <ExternalLink size="1rem" className="mr-2" />
                Browse
              </Link>
            </Button>
          </div>
        </form>
      </main>
      <footer>
        <h2 className="text-2xl font-semibold mb-4">Resolved Disputes</h2>
        <div className="flex flex-col md:grid md:grid-cols-3 gap-4">
          {data!.archives.map((props, i) => (
            <ArchivedDispute key={i} {...props} />
          ))}
        </div>
      </footer>
    </div>
  );
}
