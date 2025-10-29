import { ReactNode } from "react";
import { create } from "zustand";

export type UseConfirmStore = {
  confirm: (() => void) | null;
  confirmContent: ReactNode;
  setConfirm: (confirm: (() => void) | null) => void;
  setConfirmContent: (content: ReactNode) => void;
};

export const useConfirmStore = create<UseConfirmStore>()((set) => ({
  confirm: null,
  confirmContent: null,
  setConfirm: (confirm) => set({ confirm }),
  setConfirmContent: (content) => set({ confirmContent: content }),
}));
