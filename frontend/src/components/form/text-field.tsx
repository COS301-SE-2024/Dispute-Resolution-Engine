import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "../ui/form";
import { Input } from "../ui/input";

export default function TextField<T>({
  name,
  label,
  type,
}: {
  name: keyof T;
  label: string;
  type?: "text" | "password";
}) {
  return (
    <FormField
      name={name.toString()}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{label}</FormLabel>
          <FormControl>
            <Input type={type} placeholder={label} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}
