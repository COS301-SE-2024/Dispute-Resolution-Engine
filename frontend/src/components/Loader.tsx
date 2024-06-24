import Image from "next/image";

export default function Loader() {
  return (
    <div className="flex justify-center items-center h-screen">
      <div className="animate-spin bg-transparent rounded-full h-32 w-32 border-2 border-dre-bg-dark/10 dark:border-dre-300 overflow-hidden">
        <Image src="/logo.svg" width={1362} height={1362} alt="Ralph" className="h-32 w-32 mt-6" />
      </div>
    </div>
  );
}
