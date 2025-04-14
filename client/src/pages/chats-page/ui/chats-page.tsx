import { useChatsQuery } from "@/entities/chats";
import { useProfileQuery } from "@/entities/user";
import { BottomNavigation } from "@/features/navigation";
import { NAVIGATION } from "@/shared/navigation";
import { cn } from "@/shared/utils/cn";
import { PlusCircleIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { useState } from "react";
import { Link } from "wouter";
import { ChatPreview } from "./chat-preview";
import { ChatsSearch } from "./chats-search";

export const ChatsPage = () => {
  const [search, setSearch] = useState("");
  const { data, isPending, isError, error } = useChatsQuery(search);
  const profileQuery = useProfileQuery();

  if (isPending || profileQuery.isPending) return <>Loading...</>;

  if (isError || profileQuery.isError) return <>Error: {error?.message}</>;

  const handleSearchSubmit = (searchValue: string) => {
    setSearch(searchValue);
  };

  return (
    <>
      <section className="p-4 pt-14">
        <div className="flex gap-2 items-center">
          <ChatsSearch onSubmitSearch={handleSearchSubmit} className="grow" />

          <Button
            as={Link}
            isIconOnly
            color="primary"
            href={NAVIGATION.createChat}
            state={{ origin: NAVIGATION.chats() }}
          >
            <PlusCircleIcon className="w-6 h-6" />
          </Button>
        </div>

        <ul className="mt-2">
          {data?.data.map((chat, i) => {
            const isFirstChat = i === 0;

            return (
              <li
                key={chat.id}
                className={cn({ "border-t border-t-black/20": !isFirstChat })}
              >
                <Link href={NAVIGATION.chats(chat.id)}>
                  <ChatPreview
                    key={chat.id}
                    chatPreview={chat}
                    currentUser={profileQuery.data}
                  />
                </Link>
              </li>
            );
          })}
        </ul>
      </section>

      <BottomNavigation className="fixed bottom-4 left-1/2 -translate-x-1/2" />
    </>
  );
};
