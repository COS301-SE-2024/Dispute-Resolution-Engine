import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { SelectProps } from "@radix-ui/react-select";

const GENDERS = ["Male", "Female", "Non-binary", "Prefer not to say", "Other"];

interface GenderSelectProps extends SelectProps {}

export default function GenderSelect(props: GenderSelectProps) {
  return (
    <Select {...props}>
      <SelectTrigger>
        <SelectValue placeholder="Select a gender" />
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
