import { IMessage } from "@/entities/chats";

export type IResponse<T> = { data: T };

export type IUser = {
  id: string;
  username: string;
  name: string;
};
export type IChat = { id: string; users: IUser[]; messages: IMessage[] };
export type IChatPreview = Pick<IChat, "id" | "users" | "messages">;
