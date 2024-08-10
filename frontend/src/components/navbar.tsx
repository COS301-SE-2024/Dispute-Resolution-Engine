import { CircleUserRound } from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { Button, buttonVariants } from "./ui/button";
import { cookies } from "next/headers";

const link = buttonVariants({ variant: "link" });

export default function Navbar() {
  const result = cookies().get("jwt");
  return (
    <nav className="bg-dre-200 px-4 py-2 flex items-center fixed w-full z-50">
      <Link href="/splash">
        <Image src="/logo.svg" alt="DRE Logo" width={64} height={64} />
      </Link>

      <div className="grow">
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
        {result && (
          <Link className={link} href="/workflow">
            Workflows
          </Link>
        )}
      </div>
      {result ? (
        <Link href="/profile">
          <CircleUserRound />
        </Link>
      ) : (
        <Button asChild variant="link">
          <Link href="/signup">Signup</Link>
        </Button>
      )}
    </nav>
  );
}
