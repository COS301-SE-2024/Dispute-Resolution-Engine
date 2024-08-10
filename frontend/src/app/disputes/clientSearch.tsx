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
    <div>
      <Input
        placeholder="Search"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />
      <nav className="h-full">
        <Suspense fallback={<Loader />}>
          <ul>
            {filteredData.length > 0 ? (
              filteredData.map((d) => (
                <li key={d.id}>
                  <DisputeLink href={`/disputes/${d.id}`}>
                    {d.title}
                    {d.role == "Complainant" ? (
                      <Badge className="ml-2">{d.role.substring(0, 1)}</Badge>
                    ) : d.role == "Respondent" ? (
                      <Badge className="ml-2" variant="secondary">
                        {d.role.substring(0, 1)}
                      </Badge>
                    ) : null}
                  </DisputeLink>
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
    </div>
  );
}

