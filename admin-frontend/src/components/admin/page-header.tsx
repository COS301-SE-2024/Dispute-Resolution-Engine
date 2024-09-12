import { ReactNode } from 'react';

export default function PageHeader({ label, children }: { label: string; children?: ReactNode }) {
  return (
    <header className="flex">
      <h2 className="grow p-5 font-bold tracking-wide text-xl border-b dark:border-primary-500/30 border-primary-500/20">
        {label}
      </h2>
      {children && <div>{children}</div>}
    </header>
  );
}
