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

type LanguageSelectProps = SelectProps & {
  id: string;
};

export default async function LanguageSelect(props: LanguageSelectProps) {
  const data = (await fetchLanguages()).data!;

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
