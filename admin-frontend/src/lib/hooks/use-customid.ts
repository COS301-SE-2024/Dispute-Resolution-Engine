import { useRef } from "react";

/**
 * Used for assigning IDs to both nodes an edges. This is required because
 * useId cannot be called inside a useCallback function, so a custom
 * implementation is required.
 */
export function useCustomId(start: number | undefined) {
  let count = useRef(start ?? 0);
  return function () {
    const id = count.current.toString();
    count.current++;
    return id;
  };
}
