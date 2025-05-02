import { ISearchUsersResponse, IUser } from "@/entities/user";
import { ExclamationCircleIcon } from "@heroicons/react/24/outline";
import { motion } from "framer-motion";
import { ReactNode } from "react";
import { CREATE_CHAT_PAGE_ANIMATIONS } from "../constants/animations";

type ChildrenOptions = {
  user: IUser;
  isNeedRenderPaginationTrigger: boolean;
  isSelected: boolean;
};

type SearchUsersListProps = {
  users: ISearchUsersResponse["data"];
  selectedUsers: IUser[];
  hasNextUsersPage: boolean;
  error?: Error | null;
  children: (options: ChildrenOptions) => ReactNode;
};

export const SearchUsersList = ({
  users,
  selectedUsers,
  hasNextUsersPage,
  error,
  children,
}: SearchUsersListProps) => {
  if (error) {
    return (
      <motion.div
        className="text-large flex items-center justify-center gap-1 mt-4"
        {...CREATE_CHAT_PAGE_ANIMATIONS.FETCH_USERS_ERROR_MESSAGE}
      >
        <ExclamationCircleIcon className="w-5 h-5" />
        Oops, Something went wrong...
      </motion.div>
    );
  }

  return (
    <ul className="mt-2 flex flex-col">
      {users.map((user, i) => {
        const isPlaceForPaginationTrigger =
          i === Math.max(0, users.length - 10);

        const isNeedRenderPaginationTrigger =
          isPlaceForPaginationTrigger && hasNextUsersPage;

        const isSelected =
          selectedUsers.find((selectedUser) => selectedUser.id === user.id) !==
          undefined;

        return children?.({ user, isNeedRenderPaginationTrigger, isSelected });
      })}
    </ul>
  );
};
