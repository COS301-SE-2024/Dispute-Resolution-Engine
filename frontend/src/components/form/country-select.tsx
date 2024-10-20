"use client";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { getCountries } from "@/lib/api";
import { SelectProps } from "@radix-ui/react-select";
import { useQuery } from "@tanstack/react-query";

export default function CountrySelect({
  id,
  ...props
}: SelectProps & {
  id?: string;
}) {
  const query = useQuery({
    queryKey: ["countryList"],
    queryFn: () => getCountries(),
    staleTime: Infinity,
  });

  return (
    <Select {...props}>
      <SelectTrigger disabled={query.isPending} id={id}>
        <SelectValue placeholder="Select a country" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {query.data?.map((country) => (
            <SelectItem key={country.country_code} value={country.country_code}>
              {country.country_name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
