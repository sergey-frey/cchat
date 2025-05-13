import { IMessage } from "@/entities/chats";

export type IResponse<T> = { status: number; data: T };

export type IResponseWithCursor<T> = IResponse<
  {
    cursors: {
      prev_cursor: number | null;
      next_cursor: number | null;
    };
  } & T
>;

export type IUserServerData = {
  ID: string;
  Username: string;
  Email: string;
  Name: string;
};

export type IUser = {
  id: string;
  username: string;
  email: string;
  name: string;
};

export type IChat = { id: string; users: IUser[]; messages: IMessage[] };
export type IChatPreview = Pick<IChat, "id" | "users" | "messages">;
