import { useEffect } from "react";
import { useToast } from "./use-toast";

/**
 * Automatically displays an error toast when the passed-in error becomes non-null.
 * @param err The error to monitor for displaying the toast.
 * @param title The title of the toast.
 */
export function useErrorToast(err: Error | null, title: string) {
  const { toast } = useToast();
  useEffect(() => {
    if (err) {
      toast({
        variant: "error",
        title,
        description: err?.message,
      });
    }
  }, [err, toast]);
}
