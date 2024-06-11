import Link from "next/link"
import { Button } from "@/components/ui/button"
import Image from "next/image"

export default function SplashHeader() {
  const buttonClass : string = "group inline-flex h-9 w-max items-center justify-center rounded-md bg-white px-4 py-2 text-sm font-medium transition-colors hover:bg-gray-100 hover:text-gray-900 focus:bg-gray-100 focus:text-gray-900 focus:outline-none disabled:pointer-events-none disabled:opacity-50 data-[active]:bg-gray-100/50 data-[state=open]:bg-gray-100/50 dark:bg-gray-950 dark:hover:bg-gray-800 dark:hover:text-gray-50 dark:focus:bg-gray-800 dark:focus:text-gray-50 dark:data-[active]:bg-gray-800/50 dark:data-[state=open]:bg-gray-800/50"
  return (
    <div className="container mx-auto px-4 md:px-6 lg:px-8 bg-black z-40">
      <header className="flex h-20 w-full shrink-0 items-center px-4 md:px-6">
        <Link href="#" className="mr-6 hidden lg:flex" prefetch={false}>
          <RalphIcon className="h-6 w-6" />
          <span className="sr-only">Car E-commerce</span>
        </Link>
        <div className="ml-auto flex gap-2">
          <Link
            href="#"
            className={buttonClass}
            prefetch={false}
          >
            Home
          </Link>
          <Link
            href="#"
            className={buttonClass}
            prefetch={false}
          >
            Disputes
          </Link>
          <Button variant="outline" className="justify-self-end px-2 py-1 text-xs">
            Login
          </Button>
          <Button className="justify-self-end px-2 py-1 text-xs">
            Sign Up
          </Button>
        </div>
      </header>
    </div>
  );
}

function RalphIcon(props: any) {
  const imgSize : number = 55
  return (
    <Image
      src="/ralph.png"
      alt="Racoon Logo"
      width={imgSize}
      height={imgSize}
    />
    // <svg
    //   {...props}
    //   xmlns="http://www.w3.org/2000/svg"
    //   width="24"
    //   height="24"
    //   viewBox="0 0 24 24"
    //   fill="none"
    //   stroke="currentColor"
    //   strokeWidth="2"
    //   strokeLinecap="round"
    //   strokeLinejoin="round"
    // >
    //   <path d="M19 17h2c.6 0 1-.4 1-1v-3c0-.9-.7-1.7-1.5-1.9C18.7 10.6 16 10 16 10s-1.3-1.4-2.2-2.3c-.5-.4-1.1-.7-1.8-.7H5c-.6 0-1.1.4-1.4.9l-1.4 2.9A3.7 3.7 0 0 0 2 12v4c0 .6.4 1 1 1h2" />
    //   <circle cx="7" cy="17" r="2" />
    //   <path d="M9 17h6" />
    //   <circle cx="17" cy="17" r="2" />
    // </svg>
  )
}