import { useEffect, useState } from "react";
import { containerRefSelector, useAppContainer } from "./app-container-store";

type UseAppContainerScrollOptions = {
  onScroll?: (container: HTMLElement) => void;
};

export const useAppContainerScroll = ({
  onScroll,
}: UseAppContainerScrollOptions = {}) => {
  const containerRef = useAppContainer(containerRefSelector);
  const [scroll, setScroll] = useState(0);

  useEffect(() => {
    const container = containerRef?.current;

    if (!container) return;

    const handleScroll = () => {
      onScroll?.(container);
      setScroll(container.scrollTop);
    };

    container.addEventListener("scroll", handleScroll);

    return () => {
      container.removeEventListener("scroll", handleScroll);
    };
  }, [onScroll, containerRef]);

  return { scroll };
};
