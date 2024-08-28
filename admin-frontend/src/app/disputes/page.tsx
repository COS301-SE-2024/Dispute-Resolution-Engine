import { Search, Filter } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function Disputes() {
  return (
    <div className="flex flex-col">
      <header>
        <h2 className="p-5 font-bold tracking-wide text-xl border-b dark:border-primary-500/30 border-primary-500/20">
          Disputes
        </h2>
        <div class="flex items-center px-5 gap-2 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
          <div className="grid grid-cols-[auto_1fr] items-center grow">
            <input
              type="search"
              className="col-span-2 p-5 bg-transparent  col-start-1 row-start-1 pl-12"
              placeholder="Search disputes..."
            />
            <div className="p-5 row-start-1 col-start-1 pointer-events-none">
              <Search size={20} />
            </div>
          </div>
          <Button variant="ghost" className="gap-2">
            <Filter />
            <span>Filter by</span>
          </Button>
        </div>
      </header>
      <main className="overflow-auto p-5 grow">
        <h1>Insert table here</h1>
      </main>
    </div>
  );
}
