import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";

export default function DisputeRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex items-stretch h-full lg:w-3/4 mx-auto">
      <div className="w-56 flex flex-col p-2 gap-4">
        <Input placeholder="Search" />
        <nav>
          <ul>
            <li>
              <Button variant="link" asChild>
                <Link href="/disputes/1">Hello Sur!</Link>
              </Button>
            </li>
            <li>
              <Button variant="link" asChild>
                <Link href="/disputes/2">You stole my lunch</Link>
              </Button>
            </li>
            <li>
              <Button variant="link" asChild>
                <Link href="/disputes/3">You killed my family</Link>
              </Button>
            </li>
          </ul>
        </nav>

        <Button variant="secondary" className="mt-auto" asChild>
          <Link href="/disputes/create" className="w-full">
            + Create
          </Link>
        </Button>
      </div>
      <Separator orientation="vertical" />
      {children}
    </div>
  );
}
