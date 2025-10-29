import { useCallback } from "react";

type UseIntersectionOptions = {
  onIntersect: () => void;
};

export const useIntersection = <T extends HTMLElement>({
  onIntersect,
}: UseIntersectionOptions) => {
  return useCallback((element: T | null) => {
    const observer = new IntersectionObserver((entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          onIntersect();
        }
      }
    });

    if (!element) {
      return;
    }

    observer.observe(element);

    return () => {
      observer.disconnect();
    };
  }, []);
};
