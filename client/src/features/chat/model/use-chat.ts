import { getChatSocket } from "@/shared/api/socket";
import { useCallback, useEffect, useState } from "react";
import { IMessageWithAuthor, SendMessageDto } from "../types";

type UseChatOptions = {
  chatId: string;
  onSendMessage?: (message: SendMessageDto) => void;
  onReceiveMessages?: (messages: IMessageWithAuthor[]) => void;
};

export const useChat = ({
  chatId,
  onSendMessage,
  onReceiveMessages,
}: UseChatOptions) => {
  const [messages, setMessages] = useState<IMessageWithAuthor[]>([]);

  const [socket, setSocket] = useState(() => getChatSocket(chatId));

  const sendMessage = useCallback(
    (message: SendMessageDto) => {
      return new Promise<void>((resolve, reject) => {
        if (!message.content.trim()) {
          resolve();
          return;
        }

        try {
          socket.send(JSON.stringify(message));
          onSendMessage?.(message);
          resolve();
        } catch (error) {
          if (socket.CLOSING || socket.CLOSED) {
            setSocket(getChatSocket(chatId));
          }

          reject(error);
        }
      });
    },
    [socket, chatId, onSendMessage],
  );

  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      try {
        const messages = JSON.parse(event.data) as IMessageWithAuthor[];
        setMessages((prev) => [...prev, ...messages]);
        onReceiveMessages?.(messages);
      } catch (error) {
        console.error(error);
      }
    };

    socket.addEventListener("message", handleMessage);

    return () => {
      socket.removeEventListener("message", handleMessage);
    };
  }, [socket, onReceiveMessages]);

  return { messages, methods: { sendMessage } };
};
