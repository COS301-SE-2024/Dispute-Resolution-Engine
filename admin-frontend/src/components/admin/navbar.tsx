"use client";

import {
  LayoutDashboard,
  TicketCheck,
  Network,
  FileText,
  Menu,
  Settings,
  HelpCircle,
  Users,
} from "lucide-react";
import { useState, useMemo } from "react";

import Ralph from "@/components/ralph";
import NavLink from "./nav-link";
import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";

function SectionTitle({ label, show }: { label: string; show: boolean }) {
  return (
    <p
      className={cn(
        "text-black/75 dark:text-white/75 tracking-wide mb-1 transition-opacity",
        "overflow-x-visible w-0 text-nowrap",
        !show && "opacity-0 select-none "
      )}
    >
      {label}
    </p>
  );
}

export default function Navbar({ className }: { className?: string }) {
  const [expanded, setExpanded] = useState(true);

  const fullClass = useMemo(() => {
    return cn(
      "overflow-x-hidden  border-b md:border-b-0 md:h-full p-2 grid grid-rows-[auto_1fr_auto] md:border-r dark:border-primary-500/30 border-primary-500/20",
      "md:w-16 transition-all",
      expanded && "md:w-56 shadow-lg md:shadow-none",
      className
    );
  }, [expanded, className]);

  return (
    <aside className={fullClass}>
      <header className="flex items-center gap-4 min-h-16 md:mb-5">
        <div className={cn("flex items-center gap-4 grow", !expanded && "md:hidden")}>
          <Ralph className="ml-2" />
          <h1 className="text-base tracking-widest font-bold">DRE</h1>
        </div>
        <Button
          className="p-3 h-auto rounded-full"
          variant="ghost"
          onClick={() => setExpanded(!expanded)}
        >
          <Menu />
        </Button>
      </header>
      <nav className={cn(!expanded && "hidden", "md:block")}>
        <SectionTitle label="Main Menu" show={expanded} />
        <ul className="flex flex-col gap-2">
          <li>
            <NavLink label="Dashboard" href="/" icon={<LayoutDashboard />} expanded={expanded} />
          </li>
          <li>
            <NavLink label="Tickets" href="/tickets" icon={<TicketCheck />} expanded={expanded} />
          </li>
          <li>
            <NavLink label="Workflows" href="/workflows" icon={<Network />} expanded={expanded} />
          </li>
          <li>
            <NavLink label="Disputes" href="/disputes" icon={<FileText />} expanded={expanded} />
          </li>
          <li>
            <NavLink label="Experts" href="/experts" icon={<Users />} expanded={expanded} />
          </li>
        </ul>
      </nav>
      <footer className={cn(!expanded && "hidden", "md:block")}>
        <SectionTitle label="Other" show={expanded} />
        <ul>
          <li>
            <NavLink label="Settings" href="/settings" icon={<Settings />} expanded={expanded} />
          </li>
          <li>
            <NavLink label="Help" href="/help" icon={<HelpCircle />} expanded={expanded} />
          </li>
        </ul>
      </footer>
    </aside>
  );
}
