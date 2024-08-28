"use client";

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectItem,
} from "@/components/ui/select";
import { SearchParams } from "./page";
import { useState } from "react";
import { Button } from "@/components/ui/button";

export default function SearchFilters({ params }: { params: SearchParams }) {
  const [changed, setChanged] = useState(false);

  return (
    <div className="inline-flex ml-3 gap-3">
      <Select name="sort" defaultValue={params.sort} onValueChange={() => setChanged(true)}>
        <SelectTrigger className="min-w-40">
          <SelectValue placeholder="Sort by" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem value="title">Title</SelectItem>
            <SelectItem value="date_filed">Date filed</SelectItem>
            <SelectItem value="date_resolved">Date resolved</SelectItem>
            <SelectItem value="time_taken">Time taken</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
      <Select name="order" defaultValue={params.order} onValueChange={() => setChanged(true)}>
        <SelectTrigger className="min-w-40">
          <SelectValue placeholder="Order" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem value="asc">Ascending</SelectItem>
            <SelectItem value="desc">Descending</SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
      {changed && <Button type="submit">Update</Button>}
    </div>
  );
}
