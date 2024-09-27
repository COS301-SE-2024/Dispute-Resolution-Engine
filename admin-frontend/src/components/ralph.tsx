"use client";

import Link from "next/link";
import Image from "next/image";

import { cn } from "@/lib/utils";

export default function Ralph({ className, href }: { className?: string; href?: string }) {
  const classes = cn("w-[50px] h-[50px] relative flex items-center justify-center", className);
  const children = (
    <>
      <div
        className="absolute w-[30px] h-[30px] rounded-full blur-[15px] bg-primary-400"
        aria-hidden="true"
      />
      <Image src="/logo.svg" alt="DRE Logo" width={50} height={50} className="z-10" />
    </>
  );

  return href ? (
    <Link href={href} className={classes}>
      {children}
    </Link>
  ) : (
    <figure className={classes}>{children}</figure>
  );
}
