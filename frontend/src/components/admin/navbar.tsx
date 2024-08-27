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
        "text-white/75 tracking-wide mb-1 transition-opacity",
        "overflow-x-visible w-0 text-nowrap",
        !show && "opacity-0 select-none ",
      )}
    >
      {label}
    </p>
  );
}

export default function Navbar({ className }: { className?: string }) {
  const [expanded, setExpanded] = useState();

  const fullClass = useMemo(() => {
    return cn(
      "overflow-x-hidden  border-b md:h-full p-2 grid grid-rows-[auto_1fr_auto] md:border-r-2 border-primary-500/30",
      expanded ? "md:w-56 shadow-lg" : "md:w-fit",
      className,
    );
  }, [expanded]);

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
            <NavLink
              label="Dashboard"
              href="/admin"
              icon={<LayoutDashboard />}
              expanded={expanded}
            />
          </li>
          <li>
            <NavLink
              label="Tickets"
              href="/admin/tickets"
              icon={<TicketCheck />}
              expanded={expanded}
            />
          </li>
          <li>
            <NavLink
              label="Workflows"
              href="/admin/workflows"
              icon={<Network />}
              expanded={expanded}
            />
          </li>
          <li>
            <NavLink
              label="Disputes"
              href="/admin/disputes"
              icon={<FileText />}
              expanded={expanded}
            />
          </li>
          <li>
            <NavLink label="Experts" href="/admin/experts" icon={<Users />} expanded={expanded} />
          </li>
        </ul>
      </nav>
      <footer className={cn(!expanded && "hidden", "md:block")}>
        <SectionTitle label="Other" show={expanded} />
        <ul>
          <li>
            <NavLink
              label="Settings"
              href="/admin/settings"
              icon={<Settings />}
              expanded={expanded}
            />
          </li>
          <li>
            <NavLink label="Help" href="/admin/help" icon={<HelpCircle />} expanded={expanded} />
          </li>
        </ul>
      </footer>
    </aside>
  );
}
