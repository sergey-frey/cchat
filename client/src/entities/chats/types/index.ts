export enum MessageType {
  TEXT = "TEXT",
}

export type IMessage = {
  id: string;
  type: MessageType;
  content: string;
  date: Date;
};
