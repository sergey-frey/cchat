// Dependencies
import { ArrowRightIcon, ChevronLeftIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { motion } from "framer-motion";
import { useState } from "react";
import { Link } from "wouter";

import { IUser } from "@/entities/user";
import { NavigationOriginState } from "@/shared/types/navigation";
import { cn } from "@/shared/utils/cn";
import { useArrayState } from "@/shared/utils/use-array-state";
import { useLocationState } from "@/shared/utils/use-location-state";
import { CREATE_CHAT_PAGE_ANIMATIONS } from "../constants/animations";
import { useCreateChat } from "../model/use-create-chat";
import { useCreateChatHandlers } from "../model/use-create-chat-handlers";
import { ChatMemberBadge } from "./chat-member-badge";
import { ChatMembers } from "./chat-members";
import { CreateChatUserItem } from "./create-chat-user-item";
import { ScrollTopButton } from "./scroll-top-button";
import { SearchUsersList } from "./search-users-list";
import { CREATE_CHAT_SEARCH_LIMIT } from "../constants";

export const CreateChatPage = () => {
  const { state } = useLocationState<NavigationOriginState>();
  const [search, setSearch] = useState("");
  const [chatMembers, chatMembersMethods] = useArrayState<IUser>([]);

  const {
    users,
    searchInputEndContent,
    hasNextUsersPage,
    isShowCreateChatButton,
    fetchUsersError,
    isShowPlaceholders,
    isShowScrollDivider,
    isChatMembersDirty,
    paginationTriggerRef,
    debouncedRefetchUsers,
  } = useCreateChat({ search, chatMembers });

  const { handleScrollTopClick, handleSearchChange, handleUserClick } =
    useCreateChatHandlers({
      setSearch,
      debouncedRefetchUsers,
      chatMembers,
      chatMembersMethods,
    });

  return (
    <section className="p-4 pt-0 relative">
      <motion.div
        className={cn(
          "pt-4 pb-4 sticky top-0 z-10",
          "bg-white border-b border-transparent transition-colors",
          isShowScrollDivider && "border-foreground-300",
          isChatMembersDirty && "pb-2 border-foreground-300",
        )}
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
        queryLimit={CREATE_CHAT_SEARCH_LIMIT}
        selectedUsers={chatMembers}
        error={fetchUsersError}
      >
        {({ user, isNeedRenderPaginationTrigger, isSelected }) => {
          return (
            <CreateChatUserItem
              key={user.id}
              user={user}
              isSelected={isSelected}
              isShowPlaceholders={isShowPlaceholders}
              onClick={handleUserClick}
              endContent={
                isNeedRenderPaginationTrigger && (
                  <div ref={paginationTriggerRef} />
                )
              }
            />
          );
        }}
      </SearchUsersList>

      <ScrollTopButton onPress={handleScrollTopClick} className="z-[101]" />
    </section>
  );
};
