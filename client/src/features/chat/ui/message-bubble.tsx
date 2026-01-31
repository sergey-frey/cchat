import { HTMLAttributes } from "react";
import { IMessageWithAuthor } from "../types";
import { cn } from "@/shared/utils/cn";

type MessageBubbleProps = HTMLAttributes<HTMLSpanElement> & {
  message: IMessageWithAuthor;
  isOwn: boolean;
  isBlockStart?: boolean;
  isBlockEnd?: boolean;
};

export const MessageBubble = ({
  message,
  isOwn,
  isBlockStart,
  isBlockEnd,
  className,
  ...props
}: MessageBubbleProps) => {
  return (
    <span
      {...props}
      className={cn(
        "px-3 py-1.5 rounded-md bg-slate-500 text-gray-100 w-fit",
        "wrap-break-word overflow-hidden",
        isOwn && "bg-primary-400 text-white",
        isBlockStart && "rounded-t-xl",
        isBlockEnd && "rounded-b-xl",
        className,
      )}
    >
      {message.content}
    </span>
  );
};
