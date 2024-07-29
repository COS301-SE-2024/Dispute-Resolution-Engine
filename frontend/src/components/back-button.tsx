"use client";

import { useRouter } from "next/navigation";
import { Button, ButtonProps } from "./ui/button";

export type BackButtonProps = Pick<ButtonProps, Exclude<keyof ButtonProps, "onClick">>;

export function BackButton(props: BackButtonProps) {
  const router = useRouter();

  return <Button onClick={() => router.back()} {...props} />;
}
