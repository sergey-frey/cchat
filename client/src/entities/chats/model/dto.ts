import { IUser } from "@/shared/api/types";

export type CreateChatDto = {
  users: IUser["id"][];
};
