import { IUser, UserBadge } from "@/entities/user";
import { TrashIcon, XCircleIcon } from "@heroicons/react/24/outline";
import { cn } from "@heroui/theme";
import { motion } from "framer-motion";
import { twJoin } from "tailwind-merge";
import { CREATE_CHAT_PAGE_ANIMATIONS } from "../constants/animations";

type ChatMemberBadgeProps = {
  user: IUser;
  onClick: () => void;
};

export const ChatMemberBadge = ({ user, onClick }: ChatMemberBadgeProps) => {
  return (
    <motion.button
      className="inline"
      onClick={onClick}
      {...CREATE_CHAT_PAGE_ANIMATIONS.CHAT_MEMBER_BADGE}
    >
      <UserBadge
        className={twJoin(
          "group transition-colors",
          "hover:bg-red-200 hover:border-red-400",
        )}
      >
        <UserBadge.Avatar
          className={cn("relative transition-colors", "group-hover:bg-red-400")}
        >
          <TrashIcon
            className={cn(
              "abs-center w-4 h-4",
              "scale-50 opacity-0 transition-all",
              "group-hover:scale-100 group-hover:opacity-100",
            )}
          />
        </UserBadge.Avatar>
        <UserBadge.Username className="transition-colors group-hover:text-black">
          @{user.username}
        </UserBadge.Username>
        <UserBadge.EndContent className="transition-colors text-foreground-300 lg:hidden group-hover:text-foreground-900">
          <XCircleIcon className="w-5 h-5" />
        </UserBadge.EndContent>
      </UserBadge>
    </motion.button>
  );
};
