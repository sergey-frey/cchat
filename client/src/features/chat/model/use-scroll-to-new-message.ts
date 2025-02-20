import { useLayoutEffect, useRef } from "react";

type UseScrollToNewMessageOptions = {
  isEnabled?: boolean;
};

export const useScrollToNewMessage = ({
  isEnabled,
}: UseScrollToNewMessageOptions = {}) => {
  const rowRef = useRef<HTMLDivElement>(null);

  useLayoutEffect(() => {
    if (!isEnabled) return;

    const row = rowRef.current;

    if (!row) return;

    const parent = row.parentNode;

    if (!parent) return;

    if (parent.lastElementChild !== row) return;

    parent.lastElementChild?.scrollIntoView({ behavior: "smooth" });
  }, [isEnabled]);

  return rowRef;
};
