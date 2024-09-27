"use client";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { fetchLanguages } from "@/lib/api";
import { SelectProps } from "@radix-ui/react-select";
import { useQuery } from "@tanstack/react-query";

type LanguageSelectProps = SelectProps & {
  id: string;
};

export default function LanguageSelect({ id, ...props }: LanguageSelectProps) {
  const query = useQuery({
    queryKey: ["languageList"],
    queryFn: () => fetchLanguages(),
    staleTime: Infinity,
  });

  return (
    <Select {...props}>
      <SelectTrigger disabled={query.isPending}>
        <SelectValue id={id} placeholder="Select a language" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {query.data?.map((country) => (
            <SelectItem key={country.id} value={country.id}>
              {country.label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
