"use client";
import { getDisputeList } from "@/lib/api/dispute";
import { DisputeLink } from "./link";
import { Input } from "@/components/ui/input";
import { Suspense, useEffect, useMemo, useState } from "react";
import Loader from "@/components/Loader";
import { DisputeListResponse } from "@/lib/interfaces/dispute";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function ClientSearch() {
  const baseDLR = useMemo<DisputeListResponse>(() => [], []);
  const [data, setData] = useState(baseDLR);
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      const result = await getDisputeList();
      setData(result.data ?? baseDLR);
    };
    fetchData();
  }, [baseDLR]);

  const filteredData = useMemo(() => {
    return data.filter((d) => d.title.toLowerCase().includes(searchTerm.toLowerCase()));
  }, [data, searchTerm]);

  return (
    <div className="grid grid-rows-[auto_1fr_auto] gap-3 overflow-y-hidden p-3 md:border-r border-r-primary-500/30">
      <Input
        placeholder="Search"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />
      <nav className="overflow-y-auto">
        <Suspense fallback={<Loader />}>
          <ul className="overflow-y-auto">
            {filteredData.length > 0 ? (
              filteredData.map((d) => (
                <li key={d.id} className="border-b border-primary-500/30 py-3 mx-2">
                  <DisputeLink dispute={d.id} role={d.role} title={d.title} />
                </li>
              ))
            ) : (
              <p role="alert" className="text-dre-bg-light/50 w-full">
                You aren&apos;t involved in any disputes. Yay :)
              </p>
            )}
          </ul>
        </Suspense>
      </nav>

      <Button asChild>
        <Link href="/disputes/create">+ Create</Link>
      </Button>
    </div>
  );
}
