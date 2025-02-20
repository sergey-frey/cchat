import { cn } from "@/shared/utils/cn";
import { RefProp } from "@/shared/utils/types";
import { HTMLAttributes } from "react";

type MessagesListProps = RefProp<HTMLUListElement> &
  HTMLAttributes<HTMLUListElement>;

export const MessagesList = ({
  className,
  elemRef,
  ...props
}: MessagesListProps) => {
  return (
    <ul
      className={cn("p-4", "grid gap-1", className)}
      ref={elemRef}
      {...props}
    />
  );
};
