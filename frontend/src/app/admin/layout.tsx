import type { Metadata } from "next";
import { Rubik } from "next/font/google";
import Link from "next/link";
import Image from "next/image";

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

import "@/app/globals.css";

import Navbar from "@/components/navbar";
import Ralph from "@/components/ralph";
import { Button } from "@/components/ui/button";
import NavLink from "@/components/admin/nav-link";

const inter = Rubik({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Admin",
  icons: "/logo.svg",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <aside className="w-56 h-full p-2 grid grid-rows-[auto_1fr_auto] border-r-2 border-r-primary-500/30">
          <header className="flex items-center gap-4">
            <PanelLeftOpen />
            <Ralph />
            <strong>DRE</strong>
          </header>
          <nav>
            <p>Main Menu</p>
            <ul className="flex flex-col">
              <li>
                <NavLink label="Dashboard" href="/dashboard" icon={<LayoutDashboard />} expanded />
              </li>
              <li>
                <NavLink label="Tickets" href="/tickets" icon={<TicketCheck />} expanded />
              </li>
              <li>
                <NavLink label="Workflows" href="/workflows" icon={<Network />} expanded />
              </li>
              <li>
                <NavLink label="Disputes" href="/disputes" icon={<FileText />} expanded />
              </li>
              <li>
                <NavLink label="Experts" href="/experts" icon={<Users />} expanded />
              </li>
            </ul>
          </nav>
          <footer>
            <p>Other</p>
            <ul>
              <li>
                <NavLink label="Settings" href="/settings" icon={<Settings />} expanded />
              </li>
              <li>
                <NavLink label="Help" href="/help" icon={<HelpCircle />} expanded />
              </li>
            </ul>
          </footer>
        </aside>
      </body>
    </html>
  );
}
