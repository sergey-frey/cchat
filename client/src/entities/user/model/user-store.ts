import { IUser } from "@/shared/api/types";
import { create } from "zustand";

export type UserStore = {
  user: IUser | null | undefined;
  setUser: (user: IUser | null | undefined) => void;
};

export const useUserStore = create<UserStore>((set) => ({
  user: undefined,
  setUser: (user) => set({ user }),
}));
