import { useChatInfo } from "@/features/chat";
import { IChatPreview, IUser } from "@/shared/api/types";
import { cn } from "@/shared/utils/cn";
import { getMessageDate } from "@/shared/utils/date";

type ChatPreviewProps = {
  currentUser: IUser;
  chatPreview: IChatPreview;
};

export const ChatPreview = ({ currentUser, chatPreview }: ChatPreviewProps) => {
  const { chatTitle, lastMessage } = useChatInfo({ currentUser, chatPreview });

  return (
    <article className={cn("px-3 py-4")}>
      <h2 className="text-lg">{chatTitle}</h2>

      {lastMessage && (
        <div className="flex justify-between items-center mt-1">
          <p className="text-black/40 text-sm">{lastMessage?.content}</p>
          <time className="text-black/40 text-xs">
            {getMessageDate(lastMessage?.date)}
          </time>
        </div>
      )}
    </article>
  );
};
