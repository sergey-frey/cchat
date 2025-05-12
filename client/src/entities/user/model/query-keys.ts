import { IUser } from "../types";

export const USER_QUERY_KEYS = {
  PROFILE: "profile",
  USERS_SEARCH: "users-search",
  user: (username: IUser["username"]) => `user/${username}`,
} as const;

export type UserQueyKeyType =
  (typeof USER_QUERY_KEYS)[keyof typeof USER_QUERY_KEYS];
