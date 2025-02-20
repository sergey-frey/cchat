import { AuthFormState } from "@/pages/auth-page/types";

export const NAVIGATION = {
  auth: (state: AuthFormState) => `~/auth?state=${state}`,
  profile: "~/app/profile",
  chats: (id?: string) => (id ? `~/app/chats/${id}` : "~/app/chats"),
  createChat: "~/app/chats/create",
};
