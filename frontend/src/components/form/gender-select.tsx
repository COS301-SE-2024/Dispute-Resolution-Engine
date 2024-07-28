import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { GENDERS } from "@/lib/constants";
import { SelectProps } from "@radix-ui/react-select";

type GenderSelectProps = SelectProps & {
  id: string;
};

export default function GenderSelect({ id, ...props }: GenderSelectProps) {
  return (
    <Select {...props}>
      <SelectTrigger>
        <SelectValue id={id} placeholder="Select a gender" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {GENDERS?.map((gen) => (
            <SelectItem key={gen} value={gen}>
              {gen}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
