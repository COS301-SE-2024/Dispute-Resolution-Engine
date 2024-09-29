import { FormEvent, useEffect, useRef, useState } from "react";
import { Input } from "../ui/input";
import { cn } from "@/lib/utils";

export default function WorkflowTitle({
  value,
  onValueChange = () => {},
  disabled,
}: {
  value: string;
  onValueChange?: (value: string) => void;
  disabled?: boolean;
}) {
  const [isEditing, setEditing] = useState(false);
  useEffect(() => {
    if (disabled) {
      setEditing(false);
    }
  }, [disabled]);

  const inputRef = useRef<HTMLInputElement | null>(null);
  function onCommit(event: FormEvent<HTMLFormElement> | undefined = undefined) {
    if (event) {
      event.preventDefault();
    }

    setEditing(false);
    const newValue = inputRef.current!.value.trim();
    if (newValue.length > 0 && value != newValue) {
      onValueChange(newValue);
    }
  }
  function onEdit() {
    if (!disabled) {
      setEditing(true);
    }
  }

  const classes = "p-3 font-bold tracking-wide text-xl";

  return (
    <form className={cn("p-2", disabled && "opacity-50")} onSubmit={onCommit}>
      {!isEditing ? (
        <h2 className={classes} onClick={onEdit} onFocus={onEdit}>
          {value}
        </h2>
      ) : (
        <Input
          autoFocus
          ref={inputRef}
          onBlur={() => onCommit()}
          defaultValue={value}
          className={cn(classes, "w-fit")}
        />
      )}
    </form>
  );
}
