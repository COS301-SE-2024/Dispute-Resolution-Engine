import { Search, Filter } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function Disputes() {
  return (
    <div className="flex flex-col">
      <header>
        <h2 className="p-5 font-bold tracking-wide text-xl border-b dark:border-primary-500/30 border-primary-500/20">
          Disputes
        </h2>
        <div class="flex items-center pl-5 pr-2 border-b dark:border-primary-500/30 border-primary-500/20">
          <Search className="pointer-events-none" size={20} />
          <input
            type="search"
            className="grow p-5 bg-transparent"
            placeholder="Search disputes..."
          />
          <Button variant="ghost">
            <Filter />
            <span>Filter by</span>
          </Button>
        </div>
      </header>
      <main className="overflow-auto p-5 grow">
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
        <h1>Hello</h1>
      </main>
    </div>
  );
}
