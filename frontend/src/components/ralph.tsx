import Link from "next/link";
import Image from "next/image";

export default function Ralph({ href }: { href?: string }) {
  const elem = href ? Link : "figure";

  return (
    <elem href={href} className="w-[50px] h-[50px] relative flex items-center justify-center">
      {/* TODO: Figure out a way for this to only render in dark mode */}
      <div
        className="absolute w-[30px] h-[30px] rounded-full blur-[15px] bg-primary-200"
        aria-hidden="true"
      />
      <Image src="/logo.svg" alt="DRE Logo" width={50} height={50} className="z-10" />
    </elem>
  );
}
