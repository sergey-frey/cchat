import { IUser, UserItem } from "@/entities/user";
import { cn } from "@heroui/theme";
import { ReactNode } from "react";

type CreateChatUserItemProps = {
  user: IUser;
  isSelected: boolean;
  isShowPlaceholders: boolean;
  endContent?: ReactNode;
  onClick: (user: IUser) => void;
};

export const CreateChatUserItem = ({
  user,
  isSelected,
  isShowPlaceholders,
  endContent,
  onClick,
}: CreateChatUserItemProps) => {
  return (
    <button key={user.id} className="text-start" onClick={() => onClick(user)}>
      <UserItem
        className={cn(
          "py-1.5",
          isSelected && "bg-linear-to-r from-transparent to-slate-200",
        )}
      >
        <UserItem.Avatar isLoaded={!isShowPlaceholders} />
        <UserItem.Content>
          <UserItem.Name isLoaded={!isShowPlaceholders}>
            {user.name}
          </UserItem.Name>
          <UserItem.Username className="mt-0.5" isLoaded={!isShowPlaceholders}>
            @{user.username}
          </UserItem.Username>
        </UserItem.Content>
      </UserItem>
      {endContent}
    </button>
  );
};
