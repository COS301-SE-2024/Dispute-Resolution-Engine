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

  useEffect(() => {
    const fetchData = async () => {
      const result = await getDisputeList();
      setData(result.data ?? baseDLR);
    };
  
    fetchData();
  }, [baseDLR]);
  return (
    <div>
      <Input placeholder="Search" />
      <nav>
        <Suspense fallback={<Loader />}>
          <ul>
            {data ? (
              data.map((d) => (
                <li key={d.id}>
                  <DisputeLink href={`/disputes/${d.id}`}>
                    {d.title}
                    {d.role == "Complainant" ? (
                      <Badge className="ml-2">{d.role.substring(0, 1)}</Badge>
                    ) : d.role == "Respondant" ? (
                      <Badge className="ml-2" variant="secondary">
                        {d.role.substring(0, 1)}
                      </Badge>
                    ) : null}
                  </DisputeLink>
                </li>
              ))
            ) : (
              <li>Error</li>
            )}
          </ul>
        </Suspense>
      </nav>
    </div>
  );
}
