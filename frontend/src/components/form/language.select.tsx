"use client";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Language, fetchLanguages } from "@/lib/api";
import { SelectProps } from "@radix-ui/react-select";
import { useEffect, useState } from "react";

type LanguageSelectProps = SelectProps & {
  id: string;
};

export default function LanguageSelect(props: LanguageSelectProps) {
  const [data, setData] = useState<Language[]>([]);
  useEffect(() => {
    let cancelled = false;
    async function load() {
      const data = (await fetchLanguages()).data!;
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
        <SelectValue placeholder="Select a language" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {data?.map((country) => (
            <SelectItem key={country.id} value={country.id}>
              {country.label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
