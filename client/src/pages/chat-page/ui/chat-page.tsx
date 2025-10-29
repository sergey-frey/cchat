import { MessageBubble, MessageRow, useChat } from "@/features/chat";
import { useProfileQuery } from "@/entities/user";
import { NAVIGATION } from "@/shared/navigation";
import { cn } from "@/shared/utils/cn";
import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Redirect, useParams } from "wouter";
import { MessagesList } from "./messages-list";
import { Navigation } from "./navigation";
import { SendMessageForm } from "./send-message-form";

export const ChatPage = () => {
  const profileQuery = useProfileQuery();
  const { chatId } = useParams();
  const { messages, methods } = useChat({
    chatId: chatId ?? "",
  });

  if (!chatId) return <Redirect to={NAVIGATION.chats()} />;

  return (
    <section className="flex flex-col h-full">
      <Navigation chatTitle={""} />
      <section className="grow">
        <MessagesList>
          {messages.map((message, i) => {
            const isOwn = profileQuery.data?.id === message.author.id;

            let isBlockStart = true;
            let isBlockEnd = true;

            if (i > 0) {
              const prevMessage = messages[i - 1];

              isBlockStart = prevMessage.author.id !== message.author.id;
            }

            if (i < messages.length - 1) {
              const nextMessage = messages[i + 1];

              isBlockEnd = nextMessage.author.id !== message.author.id;
            }

            return (
              <MessageRow
                key={message.id}
                message={message}
                isOwn={isOwn}
                isBlockStart={isBlockStart}
              >
                <MessageBubble
                  key={message.id}
                  message={message}
                  isOwn={isOwn}
                  isBlockStart={isBlockStart}
                  isBlockEnd={isBlockEnd}
                  className="max-w-[70%]"
                />
              </MessageRow>
            );
          })}
        </MessagesList>
      </section>

      <section
        className={cn(
          "px-4 py-3 sticky bottom-0",
          "bg-white shadow-small rounded-t-medium",
        )}
      >
        <SendMessageForm
          onSubmit={methods.sendMessage}
          submitButton={
            <Button type="submit" color="primary" variant="light" isIconOnly>
              <PaperAirplaneIcon className="w-6" />
            </Button>
          }
        />
      </section>
    </section>
  );
};
