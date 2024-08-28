import { Search } from "lucide-react";
import { forwardRef, HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

type SearchProps = Omit<HTMLAttributes<HTMLInputElement>, "type">;

const SearchInput = forwardRef<HTMLInputElement, SearchProps>(({ className, ...props }, ref) => (
  <div className="flex items-center gap-1">
    <Search className="pointer-events-none" size={20} />
    <input ref={ref} type="search" className={cn("grow bg-transparent", className)} {...props} />
  </div>
));
SearchInput.displayName = "SearchInput";

export { SearchInput };
