import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "@/app/globals.css";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "DRE Disputes",
  description: "View and manage your disputes",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className}`}>
        <div className="flex items-stretch h-full">
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
      </body>
    </html>
  );
}
