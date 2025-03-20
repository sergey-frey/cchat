import { IMessage } from "@/entities/chats";

export type IResponse<T> = { status: number; data: T };

export type IUserServerData = {
  ID: string;
  Username: string;
  Email: string;
};

export type IUser = {
  id: string;
  username: string;
  email: string;
};

export type IChat = { id: string; users: IUser[]; messages: IMessage[] };
export type IChatPreview = Pick<IChat, "id" | "users" | "messages">;
