"use client";

import { Search } from "lucide-react";
import { useEffect, useState } from "react";

export default function SearchBar({
  placeholder,
  timeout = 1000,
  onUpdate = () => {},
}: {
  placeholder: string;
  onUpdate?: (value: string) => void;
  timeout?: number;
}) {
  const [value, setValue] = useState("");

  useEffect(() => {
    const cancel = setTimeout(() => {
      onUpdate(value);
      console.log(value);
    }, timeout);
    return () => clearTimeout(cancel);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value]);

  return (
    <div className="grid grid-cols-[auto_1fr] items-center grow">
      <input
        type="search"
        className="col-span-2 p-5 bg-transparent  col-start-1 row-start-1 pl-12"
        placeholder={placeholder}
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <div className="p-5 row-start-1 col-start-1 pointer-events-none">
        <Search size={20} />
      </div>
    </div>
  );
}
