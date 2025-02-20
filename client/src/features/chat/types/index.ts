import { IMessage } from "@/entities/chats";
import { IUser } from "@/shared/api/types";

export type IMessageWithAuthor = IMessage & { author: IUser };
export type SendMessageDto = Pick<IMessage, "type" | "content">;
