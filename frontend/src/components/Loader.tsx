import Image from "next/image";

export default function Loader() {
  return (
    <div className="flex justify-center items-center h-screen">
      <div className="animate-spin bg-transparent rounded-full h-32 w-32 border-2 border-black dark:border-white/50 overflow-hidden">
        <Image src="/ralph.png" width={1362} height={1362} alt="Ralph" className="h-32 w-32 mt-6" />
      </div>
    </div>
  );
}
