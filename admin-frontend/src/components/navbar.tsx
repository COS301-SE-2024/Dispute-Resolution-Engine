import Ralph from "@/components/ralph";

export default function Navbar() {
  return (
    <header className="border-b  border-primary-500/20 dark:border-primary-500/30 px-3 py-4 flex items-center gap-5">
      <Ralph />
      <h1 className="text-lg tracking-wide font-bold">Dispute Resolution Engine</h1>
    </header>
  );
}
