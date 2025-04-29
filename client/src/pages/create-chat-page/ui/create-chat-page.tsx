import { IUser, UserItem } from "@/entities/user";
import { BottomNavigation } from "@/features/navigation";
import { NavigationOriginState } from "@/shared/types/navigation";
import { useArrayState } from "@/shared/utils/use-array-state";
import { useLocationState } from "@/shared/utils/use-location-state";
import {
  ArrowRightIcon,
  ArrowUpIcon,
  ChevronLeftIcon,
} from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { useState } from "react";
import { twJoin, twMerge } from "tailwind-merge";
import { Link } from "wouter";
import { useCreateChat } from "../model/use-create-chat";
import { ChatMemberBadge } from "./chat-member-badge";
import { ChatMembers } from "./chat-members";
import { SearchUsersList } from "./search-users-list";

export const CreateChatPage = () => {
  const { state } = useLocationState<NavigationOriginState>();
  const [search, setSearch] = useState("");
  const [chatMembers, chatMembersMethods] = useArrayState<IUser>([]);

  const {
    users,
    searchInputEndContent,
    hasNextUsersPage,
    isShowCreateChatButton,
    paginationTriggerRef,
    debouncedRefetchUsers,
  } = useCreateChat({ search, chatMembers });

  const handleSearchChange = (value: string) => {
    setSearch(value);
    debouncedRefetchUsers();
  };

  return (
    <>
      <section className="p-4 pt-0 relative">
        <div className={twJoin("pt-4 pb-2", "sticky top-0 bg-white border-b")}>
          <div className="flex gap-2">
            <Button as={Link} isIconOnly variant="flat" href={state?.origin}>
              <ChevronLeftIcon className="w-5 h-5" />
            </Button>

            <Input
              startContent={"@"}
              placeholder="Find users..."
              value={search}
              onValueChange={handleSearchChange}
              endContent={searchInputEndContent}
            />
          </div>

          <ChatMembers
            className={twJoin(isShowCreateChatButton && "mt-2")}
            users={chatMembers}
            endContent={
              <>
                {isShowCreateChatButton && (
                  <Button
                    size="sm"
                    color="primary"
                    className="rounded-full"
                    isIconOnly
                  >
                    <ArrowRightIcon className="w-5 h-5" />
                  </Button>
                )}
              </>
            }
          >
            {(user) => {
              return (
                <ChatMemberBadge
                  user={user}
                  onClick={() => chatMembersMethods.remove(user)}
                />
              );
            }}
          </ChatMembers>
        </div>

        <SearchUsersList
          users={users ?? []}
          hasNextUsersPage={hasNextUsersPage}
          selectedUsers={chatMembers}
        >
          {({ user, isNeedRenderPaginationTrigger, isSelected }) => {
            return (
              <button
                key={user.id}
                className="text-start"
                onClick={() => chatMembersMethods.pushUnique(user)}
              >
                <UserItem
                  className={twMerge(
                    "py-1.5",
                    isSelected &&
                      "bg-gradient-to-r from-transparent to-slate-200",
                  )}
                >
                  <UserItem.Avatar />
                  <UserItem.Content>
                    <UserItem.Name>{user.name}</UserItem.Name>
                    <UserItem.Username>@{user.username}</UserItem.Username>
                  </UserItem.Content>
                </UserItem>
                {isNeedRenderPaginationTrigger && (
                  <div ref={paginationTriggerRef}></div>
                )}
              </button>
            );
          }}
        </SearchUsersList>

        <Button
          className="fixed bottom-10 right-5"
          color="primary"
          variant="flat"
          isIconOnly
          onPress={() => window.scrollTo(0, 0)}
        >
          <ArrowUpIcon className="w-5 h-5" />
        </Button>
      </section>

      <BottomNavigation className="fixed bottom-4 left-1/2 -translate-x-1/2" />
    </>
  );
};
