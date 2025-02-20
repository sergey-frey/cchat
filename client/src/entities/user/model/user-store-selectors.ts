import { UserStore } from "./user-store";

export const userSelector = (state: UserStore) => state.user;

export const setUserSelector = (state: UserStore) => state.setUser;
