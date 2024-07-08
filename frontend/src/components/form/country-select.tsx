import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { fetchCountries } from "@/lib/api";
import { SelectProps } from "@radix-ui/react-select";

export default async function CountrySelect({
  id,
  ...props
}: SelectProps & {
  id?: string;
}) {
  const data = (await fetchCountries()).data!;

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
