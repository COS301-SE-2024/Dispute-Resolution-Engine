import { CircleUserRound } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { Button, buttonVariants } from "./ui/button";
import { cookies } from "next/headers";
import { cn } from "@/lib/utils";

const link = cn(buttonVariants({ variant: "link" }), "text-white");

export default function Navbar() {
  const result = cookies().get("jwt");
  return (
    <nav className="bg-dre-200 px-4 py-2 flex items-center fixed w-full z-50 ">
      <Link href="/splash">
        <Image src="/logo.svg" alt="DRE Logo" width={64} height={64} />
      </Link>

      <div className="grow">
        <Link className={link} href="https://cos301-se-2024.github.io/Dispute-Resolution-Engine/">
          Help
        </Link>
        <Link className={link} href="/">
          Home
        </Link>
        <Link className={link} href="/archive">
          Archive
        </Link>
        {result && (
          <Link className={link} href="/disputes">
            Disputes
          </Link>
        )}
      </div>
      {result ? (
        <Link href="/profile">
          <CircleUserRound />
        </Link>
      ) : (
        <div>
          <Button asChild variant="link" className="text-white">
            <Link href="/login">Login</Link>
          </Button>
          <Button asChild variant="link" className="text-white">
            <Link href="/signup">Signup</Link>
          </Button>
        </div>
      )}
    </nav>
  );
}
