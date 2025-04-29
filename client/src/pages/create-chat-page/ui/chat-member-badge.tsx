import { IUser, UserBadge } from "@/entities/user";
import { TrashIcon, XCircleIcon } from "@heroicons/react/24/outline";
import { twJoin } from "tailwind-merge";

type ChatMemberBadgeProps = {
  user: IUser;
  onClick: () => void;
};

export const ChatMemberBadge = ({ user, onClick }: ChatMemberBadgeProps) => {
  return (
    <button className="inline" onClick={onClick}>
      <UserBadge
        className={twJoin("group", "hover:bg-red-200 hover:border-red-400")}
      >
        <UserBadge.Avatar className="group-hover:hidden" />
        <UserBadge.AvatarPlaceholder
          className={twJoin(
            "hidden bg-red-400",
            "relative",
            "group-hover:inline",
          )}
        >
          <TrashIcon className="abs-center w-4 h-4" />
        </UserBadge.AvatarPlaceholder>
        <UserBadge.Username className="group-hover:text-black">
          @{user.username}
        </UserBadge.Username>
        <UserBadge.EndContent className="lg:hidden">
          <XCircleIcon className="w-5 h-5" />
        </UserBadge.EndContent>
      </UserBadge>
    </button>
  );
};
