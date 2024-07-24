'use client'
import { useEffect } from 'react'
import { Button } from "@/components/ui/button";
export default function Error(
  {
    error,
    reset,
  }:
    {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <div className="flex flex-col items-center align-middle w-full">
      <h2 className="font-bold text-center text-3xl my-10 dark:text-white text-black">Ooops something went wrong</h2>
      <Button className="w-fit mb-10"
        onClick={
          () => {reset()}
        }
      >
        Try again
      </Button>
      <p>If the error persists please contact our dev team</p>
    </div>
  )
}