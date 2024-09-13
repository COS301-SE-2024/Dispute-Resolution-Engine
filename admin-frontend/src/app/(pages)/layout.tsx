import Navbar from "@/components/admin/navbar";
import { Toaster } from "@/components/ui/toaster";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
        <div className="grid grid-cols-1 grid-rows-[auto_1fr] md:grid-rows-1 md:grid-cols-[auto_1fr] h-full overflow-hidden">
          <Navbar />
          {children}
          <Toaster />
        </div>
  );
}
