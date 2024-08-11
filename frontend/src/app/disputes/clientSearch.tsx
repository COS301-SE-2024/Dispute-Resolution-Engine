"use client";
import { getDisputeList } from "@/lib/api/dispute";
import { DisputeLink } from "./link";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Suspense, useEffect, useMemo, useState } from "react";
import Loader from "@/components/Loader";
import { DisputeListResponse } from "@/lib/interfaces/dispute";

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
    <>
      <Input
        placeholder="Search"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />
      <nav className="h-full w-60">
        <Suspense fallback={<Loader />}>
          <ul className="space-y-2">
            {filteredData.length > 0 ? (
              filteredData.map((d) => (
                <li key={d.id}>
                  <DisputeLink href={`/disputes/${d.id}`} role={d.role} title={d.title} />
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
    </>
  );
}
