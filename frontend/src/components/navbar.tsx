import { CircleUserRound } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { buttonVariants } from "./ui/button";

const link = buttonVariants({ variant: "link" });

export default function Navbar() {
  return (
    <nav className="bg-dre-200 px-4 py-2 flex items-center fixed w-full z-50">
      <Image src="/logo.svg" alt="DRE Logo" width={64} height={64} />
      <div className="grow">
        <Link className={link} href="/">
          Home
        </Link>
        <Link className={link} href="/disputes">
          Disputes
        </Link>
        <Link className={link} href="/archive">
          Archive
        </Link>
      </div>
      <Link href="/profile">
        <CircleUserRound />
      </Link>
    </nav>
  );
}
