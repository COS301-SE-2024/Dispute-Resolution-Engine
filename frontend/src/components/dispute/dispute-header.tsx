import { type ReactNode } from "react";

export default function DisputeHeader({
  id,
  label,
  startDate,
  status,
  children,
}: {
  id: string;
  label: string;
  startDate: string;
  status: string;
  children?: ReactNode | ReactNode[];
}) {
  return (
    <header className="p-4 py-6 grid grid-cols-[1fr_auto]">
      <div>
        <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight lg:text-2xl">{label}</h1>
        <p className="mb-4">Dispute ID: {id}</p>
      </div>

      <dl className="grid grid-cols-2 gap-2">
        <dt className="text-right font-bold">Started:</dt>
        <dd>{startDate}</dd>
        <dt className="text-right font-bold">Status:</dt>
        <dd>{status} (TBD)</dd>
      </dl>
      <div className="col-span-2 flex gap-2">{children}</div>
    </header>
  );
}
