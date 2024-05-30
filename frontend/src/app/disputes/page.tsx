"use client";

import * as Tabs from "@radix-ui/react-tabs";

export default function Disputes() {
  return (
    <>
      <Tabs.Root asChild>
        <main className="grid grid-cols-2">
          <Tabs.List asChild>
            <nav className="flex flex-col">
              <Tabs.Trigger role="link" value="dispute1" className="text-left">
                One
              </Tabs.Trigger>
              <Tabs.Trigger role="link" value="dispute2" className="text-left">
                Two
              </Tabs.Trigger>
              <Tabs.Trigger role="link" value="dispute3" className="text-left">
                Three
              </Tabs.Trigger>
              <Tabs.Trigger role="link" value="dispute4" className="text-left">
                Four
              </Tabs.Trigger>
            </nav>
          </Tabs.List>
          <Tabs.Content value="dispute1">
            <h1 className="text-lg">Dispute 1</h1>
          </Tabs.Content>
          <Tabs.Content value="dispute2">
            <h1 className="text-lg">Dispute 2</h1>
          </Tabs.Content>
          <Tabs.Content value="dispute3">
            <h1 className="text-lg">Dispute 3</h1>
          </Tabs.Content>
          <Tabs.Content value="dispute4">
            <h1 className="text-lg">Dispute 4</h1>
          </Tabs.Content>
        </main>
      </Tabs.Root>
    </>
  );
}
