import { ISearchUsersResponse, IUser } from "@/entities/user";
import { ReactNode } from "react";

type ChildrenOptions = {
  user: IUser;
  isNeedRenderPaginationTrigger: boolean;
  isSelected: boolean;
};

type SearchUsersListProps = {
  users: ISearchUsersResponse["data"];
  selectedUsers: IUser[];
  hasNextUsersPage: boolean;
  children: (options: ChildrenOptions) => ReactNode;
};

export const SearchUsersList = ({
  users,
  selectedUsers,
  hasNextUsersPage,
  children,
}: SearchUsersListProps) => {
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
