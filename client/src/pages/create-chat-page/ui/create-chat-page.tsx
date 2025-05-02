// Dependencies
import { ArrowRightIcon, ChevronLeftIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { motion } from "framer-motion";
import { useState } from "react";
import { Link } from "wouter";

import { IUser, UserItem } from "@/entities/user";
import { NavigationOriginState } from "@/shared/types/navigation";
import {
  containerRefSelector,
  useAppContainer,
} from "@/shared/utils/app-container";
import { cn } from "@/shared/utils/cn";
import { useArrayState } from "@/shared/utils/use-array-state";
import { useLocationState } from "@/shared/utils/use-location-state";
import { CREATE_CHAT_PAGE_ANIMATIONS } from "../constants/animations";
import { useCreateChat } from "../model/use-create-chat";
import { ChatMemberBadge } from "./chat-member-badge";
import { ChatMembers } from "./chat-members";
import { ScrollTopButton } from "./scroll-top-button";
import { SearchUsersList } from "./search-users-list";

export const CreateChatPage = () => {
  const { state } = useLocationState<NavigationOriginState>();
  const [search, setSearch] = useState("");
  const [chatMembers, chatMembersMethods] = useArrayState<IUser>([]);
  const containerRef = useAppContainer(containerRefSelector);

  const {
    users,
    searchInputEndContent,
    hasNextUsersPage,
    isShowCreateChatButton,
    fetchUsersError,
    paginationTriggerRef,
    debouncedRefetchUsers,
  } = useCreateChat({ search, chatMembers });

  const handleSearchChange = (value: string) => {
    setSearch(value);
    debouncedRefetchUsers();
  };

  const handleScrollTopClick = () => {
    if (!containerRef?.current) return;
    containerRef.current.scrollTo({ top: 0, behavior: "smooth" });
  };

  return (
    <section className="p-4 pt-0 relative">
      <motion.div
        className={cn("pt-4 pb-2 sticky top-0", "bg-white border-b")}
        {...CREATE_CHAT_PAGE_ANIMATIONS.HEADER}
      >
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
          className={cn(isShowCreateChatButton && "mt-2")}
          users={chatMembers}
          endContent={
            <>
              {isShowCreateChatButton && (
                <motion.div {...CREATE_CHAT_PAGE_ANIMATIONS.CREATE_CHAT_BUTTON}>
                  <Button
                    size="sm"
                    color="primary"
                    className="rounded-full"
                    isIconOnly
                  >
                    <ArrowRightIcon className="w-5 h-5" />
                  </Button>
                </motion.div>
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
      </motion.div>

      <SearchUsersList
        users={users ?? []}
        hasNextUsersPage={hasNextUsersPage}
        selectedUsers={chatMembers}
        error={fetchUsersError}
      >
        {({ user, isNeedRenderPaginationTrigger, isSelected }) => {
          return (
            <button
              key={user.id}
              className="text-start"
              onClick={() => chatMembersMethods.pushUnique(user)}
            >
              <UserItem
                className={cn(
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

      <ScrollTopButton onPress={handleScrollTopClick} />
    </section>
  );
};
