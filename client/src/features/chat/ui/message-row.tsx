import { cn } from "@/shared/utils/cn";
import { getMessageDate } from "@/shared/utils/date";
import { HTMLAttributes } from "react";
import { useScrollToNewMessage } from "../model/use-scroll-to-new-message";
import { IMessageWithAuthor } from "../types";

type MessageRowProps = HTMLAttributes<HTMLDivElement> & {
  message: IMessageWithAuthor;
  isOwn: boolean;
  isBlockStart?: boolean;
};

export const MessageRow = ({
  message,
  isOwn,
  isBlockStart,
  className,
  children,
  ...props
}: MessageRowProps) => {
  const rowRef = useScrollToNewMessage({ isEnabled: isOwn });
  const shouldShowAvatar = !isOwn && isBlockStart;
  const shouldRenderAvatar = !isOwn;

  return (
    <div
      {...props}
      ref={rowRef}
      className={cn(
        "flex gap-2",
        "group",
        "overflow-hidden",
        isOwn && "flex-row-reverse",
        isBlockStart && "mt-5",
        className,
      )}
    >
      {shouldRenderAvatar && (
        <span
          className={cn(
            "text-white bg-indigo-700 rounded-lg",
            "self-start grid place-items-center",
            "w-7 h-7 mt-0.5",
            !shouldShowAvatar && "opacity-0",
          )}
        >
          {message.author.username[0]}
        </span>
      )}

      {children}

      <time
        className={cn(
          "text-xs opacity-0 text-gray-400 self-end transition-opacity delay-50",
          "group-hover:opacity-100 group-active:opacity-100",
        )}
      >
        {getMessageDate(message.date)}
      </time>
    </div>
  );
};
