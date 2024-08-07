"use client";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Country, fetchCountries } from "@/lib/api";
import { SelectProps } from "@radix-ui/react-select";
import { useEffect, useState } from "react";

export default function CountrySelect({
  id,
  ...props
}: SelectProps & {
  id?: string;
}) {
  const [data, setData] = useState<Country[]>([]);
  useEffect(() => {
    let cancelled = false;
    async function load() {
      const data = (await fetchCountries()).data!;
      if (!cancelled) {
        setData(data);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <Select {...props}>
      <SelectTrigger>
        <SelectValue id={id} placeholder="Select a country" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {data?.map((country) => (
            <SelectItem key={country.country_code} value={country.country_code}>
              {country.country_name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
