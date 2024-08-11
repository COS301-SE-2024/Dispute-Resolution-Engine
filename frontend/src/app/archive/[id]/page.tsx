import { BackButton } from "@/components/back-button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { fetchArchivedDispute } from "@/lib/api/archive";
import { ArrowLeft, Sparkles } from "lucide-react";

function ErrorPage({ msg }: { msg: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-full gap-5 w-2/3 mx-auto">
      <h1 className="text-4xl font-bold tracking-wide">Oops, something went wrong :(</h1>
      <p className="dark:text-white/50">{msg}</p>
    </div>
  );
}

export default async function ArchivedPageDispute({ params }: { params: { id: string } }) {
  const { data, error } = await fetchArchivedDispute(params.id);
  if (error || !data) {
    return <ErrorPage msg={error} />;
  }

  return (
    <Card variant="page" className="md:mt-10">
      <CardHeader>
        <div className="flex items-center flex-wrap mb-3 gap-2 ">
          <BackButton
            variant="ghost"
            className="rounded-full aspect-square p-0 flex justify-center"
            aria-label="Back"
          >
            <ArrowLeft />
          </BackButton>
          <CardTitle className="grow">{data.title}</CardTitle>
          <div className="space-x-1">
            {data.category.map((cat) => (
              <Badge key={cat}>{cat}</Badge>
            ))}
          </div>
        </div>
        <dl className="grid grid-cols-2">
          <dt className="font-semibold">Date Filed:</dt>
          <dd>{data.date_filed}</dd>

          <dt className="font-semibold">Date Resolved:</dt>
          <dd>{data.date_resolved}</dd>

          <dt className="font-semibold">Decision:</dt>
          <dd>{data.resolution}</dd>
        </dl>
      </CardHeader>
      <CardContent>
        <p className="mb-5">{data.summary}</p>
        <section>
          <h4 className="text-lg font-semibold">Timeline</h4>
          <ol className="ml-3">
            {data.events.map((ev) => (
              <li key={ev.timestamp}>
                {ev.description} at {ev.timestamp}
              </li>
            ))}
          </ol>
        </section>
        <AISummary summary={data.summary} />
      </CardContent>
    </Card>
  );
}

function AISummary({ summary }: { summary?: string }) {
  return (
    <section className="border-t-2 border-gray-700">
      <h4 className="text-lg font-semibold flex items-center">
        <Sparkles size="1.2rem" className="mr-2" />
        <span>AI Summary</span>
      </h4>
      <p className="dark:text-dre-bg-light/50">{summary ?? "No summary available"}</p>
    </section>
  );
}
