import { Input } from "@/components/ui/input";
import { useRef, type FormEvent } from "react";

export default function EditForm({
  value,
  onCancel = () => {},
  onCommit = () => {},
}: {
  value: string;
  onCancel?: () => void;
  onCommit?: (value: string) => void;
}) {
  const inputRef = useRef<HTMLInputElement | null>(null);

  function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const newValue = inputRef.current!.value;
    if (newValue === value) {
      onCancel();
    } else if (newValue.length === 0) {
      onCancel();
    } else {
      onCommit(newValue);
    }
  }

  return (
    <form onSubmit={onSubmit} className="grow">
      <Input
        ref={inputRef}
        defaultValue={value}
        autoFocus
        className="w-fit"
        onBlur={onCancel}
        placeholder="Node label"
      />
    </form>
  );
}
