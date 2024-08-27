import type { Metadata } from "next";
import { Rubik } from "next/font/google";
import Link from "next/link";
import Image from "next/image";

import "@/app/globals.css";

import Ralph from "@/components/ralph";
import { Button } from "@/components/ui/button";

import NavLink from "@/components/admin/nav-link";
import Navbar from "@/components/admin/navbar";

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
        <Navbar />
      </body>
    </html>
  );
}
