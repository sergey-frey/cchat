import { RefObject } from "react";
import { create } from "zustand";

type UseAppContainer = {
  containerRef: RefObject<HTMLElement> | null;
  setContainerRef: (ref: RefObject<HTMLElement>) => void;
};

export const useAppContainer = create<UseAppContainer>((set) => ({
  containerRef: null,
  setContainerRef: (ref: RefObject<HTMLElement>) => set({ containerRef: ref }),
}));

export const containerRefSelector = (store: UseAppContainer) =>
  store.containerRef;

export const setContainerRefSelector = (store: UseAppContainer) =>
  store.setContainerRef;
