"use client";

import {
  LayoutDashboard,
  TicketCheck,
  Network,
  FileText,
  PanelLeftOpen,
  PanelLeftClose,
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
        !show && "opacity-0",
      )}
    >
      {label}
    </p>
  );
}

export default function Navbar() {
  const [expanded, setExpanded] = useState();

  const className = useMemo(() => {
    return cn(
      "transition-[width] overflow-x-hidden h-full p-2 grid grid-rows-[auto_1fr_auto] border-r-2 border-r-primary-500/30",
      expanded ? "md:w-56" : "w-fit",
    );
  }, [expanded]);

  return (
    <aside className={className}>
      <header className="flex items-center gap-4 min-h-16 mb-5">
        {expanded && (
          <>
            <Ralph />
            <strong className="grow">DRE</strong>
          </>
        )}
        <Button
          className="p-3 h-auto rounded-full"
          variant="ghost"
          onClick={() => setExpanded(!expanded)}
        >
          {expanded ? <PanelLeftClose /> : <PanelLeftOpen />}
        </Button>
      </header>
      <nav>
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
      <footer>
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
